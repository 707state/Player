package main

import (
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AlbumSingle struct {
	Title   string   `json:"title" bson:"title"`
	Artists []string `json:"artists" bson:"artists"`
	Album   string   `json:"album" bson:"album"`
}

func handleAlbumSingleGet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := NewFilterBuilder().
		WithStringField(q, "album").
		WithArrayField(q, "artists").
		WithStringField(q, "title").
		Build()
	log.Printf("Album Single endpoint called with filter: %v", filter)
	ctx := r.Context()
	cursor, err := singlesCollection.Find(ctx, filter)
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)
	var albumSingle []AlbumSingle
	if err = cursor.All(ctx, &albumSingle); err != nil {
		http.Error(w, "Failed to decode albums", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"exists": len(albumSingle) > 0})
}

func handleAlbumSinglePost(w http.ResponseWriter, r *http.Request) {
	var single AlbumSingle
	if err := json.NewDecoder(r.Body).Decode(&single); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if single.Title == "" || single.Artists == nil || len(single.Artists) == 0 || single.Album == "" {
		jsonError(w, "Title, artist and single are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	filter := bson.M{
		"title":   single.Title,
		"artists": single.Artists,
	}

	var existingAlbum Album
	err := albumCollection.FindOne(ctx, filter).Decode(&existingAlbum)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 创建新专辑并添加单曲
			newAlbum := Album{
				Title:   single.Title,
				Artists: single.Artists,
			}
			_, err = albumCollection.InsertOne(ctx, newAlbum)
			if err != nil {
				jsonError(w, "Failed to create album and add single", http.StatusInternalServerError)
				return
			}
			_, err = singlesCollection.InsertOne(ctx, single)
			if err != nil {
				jsonError(w, "Failed to insert data to collection", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SuccessResponse{Message: "Album created and single added successfully"})
			return
		}
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 检查单曲是否已存在
	single_filter := bson.M{
		"album":   single.Album,
		"artists": single.Artists,
		"title":   single.Title,
	}
	update := bson.M{
		"$setOnInsert": bson.M{
			"album":   single.Album,
			"artists": single.Artists,
			"title":   single.Title,
		},
	}
	result, err := singlesCollection.UpdateOne(
		ctx,
		single_filter,
		update,
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		jsonError(w, "Single collection update failed!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if result.UpsertedCount > 0 {
		json.NewEncoder(w).Encode(SuccessResponse{
			Message: "Single added to collection",
		})
	} else if result.ModifiedCount > 0 {
		json.NewEncoder(w).Encode(SuccessResponse{
			Message: "Single modified in collection",
		})
	} else {
		json.NewEncoder(w).Encode(SuccessResponse{
			Message: "Single already exists",
		})
	}
	return
}

func handleAlbumSingleDelete(w http.ResponseWriter, r *http.Request) {
	var single AlbumSingle
	if err := json.NewDecoder(r.Body).Decode(&single); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("AlbumSingle API endpoint called with book: %v", single)
	if single.Album == "" || single.Title == "" || len(single.Artists) == 0 {
		jsonError(w, "Album Title and Author and Cut are required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	filter := bson.M{
		"artists": single.Artists,
		"title":   single.Title,
		"album":   single.Album,
	}
	result, err := singlesCollection.DeleteOne(ctx, filter)
	if err != nil {
		jsonError(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if result.DeletedCount > 0 {
		json.NewEncoder(w).Encode(SuccessResponse{Message: "Single deleted successfully"})
	} else {
		json.NewEncoder(w).Encode(SuccessResponse{
			Message: "Single does not exists",
		})
	}
}
