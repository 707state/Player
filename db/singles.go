package db

import (
	"encoding/json"
	"log"
	"musaic/util"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AlbumSingle struct {
	Title        string    `json:"title" bson:"title"`
	Artists      []string  `json:"artists" bson:"artists"`
	Album        string    `json:"album" bson:"album"`
	LastModified time.Time `json:"last_modified" bson:"last_modified"`
}

func handleAlbumSingleGet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := util.NewFilterBuilder().
		WithStringField(q, "album").
		WithArrayField(q, "artists").
		WithStringField(q, "title").
		Build()
	ctx := r.Context()
	if _, exists := filter["title"]; exists {
		log.Printf("Album Single endpoint called with filter: %v", filter)
		cursor, err := singlesCollection.Find(ctx, filter)
		if err != nil {
			util.JsonError(w, "Database error", http.StatusInternalServerError)
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
	} else {
		log.Printf("Album Single endpoint called without title with filter: %v", filter)
		//不存在title，则返回所有单曲
		cursor, err := singlesCollection.Find(ctx, filter)
		if err != nil {
			util.JsonError(w, "Database error", http.StatusInternalServerError)
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
		json.NewEncoder(w).Encode(albumSingle)
		return
	}
}

func handleAlbumSinglePost(w http.ResponseWriter, r *http.Request) {
	var single AlbumSingle
	if err := json.NewDecoder(r.Body).Decode(&single); err != nil {
		util.JsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	single.LastModified = time.Now()

	if single.Title == "" || len(single.Artists) == 0 || single.Album == "" {
		util.JsonError(w, "Album, artist and single are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// 检查单曲是否已存在
	single_filter := bson.M{
		"album":   single.Album,
		"artists": single.Artists,
		"title":   single.Title,
	}
	update := bson.M{
		"$setOnInsert": bson.M{
			"album":         single.Album,
			"artists":       single.Artists,
			"title":         single.Title,
			"last_modified": single.LastModified,
		},
	}
	result, err := singlesCollection.UpdateOne(
		ctx,
		single_filter,
		update,
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		util.JsonError(w, "Single collection update failed!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if result.UpsertedCount > 0 {
		json.NewEncoder(w).Encode(util.SuccessResponse{
			Message: "Single added to collection",
		})
	} else if result.ModifiedCount > 0 {
		json.NewEncoder(w).Encode(util.SuccessResponse{
			Message: "Single modified in collection",
		})
	} else {
		json.NewEncoder(w).Encode(util.SuccessResponse{
			Message: "Single already exists",
		})
	}
}

func handleAlbumSingleDelete(w http.ResponseWriter, r *http.Request) {
	var single AlbumSingle
	if err := json.NewDecoder(r.Body).Decode(&single); err != nil {
		util.JsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("AlbumSingle API endpoint called with book: %v", single)
	if single.Album == "" || single.Title == "" || len(single.Artists) == 0 {
		util.JsonError(w, "Album Title and Author and Cut are required", http.StatusBadRequest)
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
		util.JsonError(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if result.DeletedCount > 0 {
		json.NewEncoder(w).Encode(util.SuccessResponse{Message: "Single deleted successfully"})
	} else {
		json.NewEncoder(w).Encode(util.SuccessResponse{
			Message: "Single does not exists",
		})
	}
}
