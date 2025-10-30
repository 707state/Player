package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Album struct {
	Title   string   `json:"title" bson:"title"`
	Artists []string `json:"artist" bson:"artist"`
	Genre   string   `json:"genre" bson:"genre"`
	Year    Date     `json:"year" bson:"year"`
	Url     string   `json:"url" bson:"url"`
	// url of artwork!
	Artwork string `json:"artwork" bson:"artwork"`
	Comment string `json:"comment" bson:"comment"`
	// between 0-5
	Rating       int       `json:"rating" bson:"rating"`
	LastModified time.Time `json:"last_modified" bson:"last_modified"`
}

func handleMusicGet(w http.ResponseWriter, r *http.Request) {

	// parse query params from URL and build a MongoDB filter
	q := r.URL.Query()
	filter := NewFilterBuilder().
		WithStringField(q, "title").
		WithArrayField(q, "artists").
		Build()
	log.Printf("Music API endpoint called with filter: %v", filter)
	ctx := r.Context()
	cursor, err := albumCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to fetch albums", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var albums []Album
	if err = cursor.All(ctx, &albums); err != nil {
		http.Error(w, "Failed to decode albums", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, "Failed to encode albums", http.StatusInternalServerError)
		return
	}
}
func handleMusicPost(w http.ResponseWriter, r *http.Request) {
	var album Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Music API endpoint called with album: %v", album)
	// 验证必需字段
	if album.Title == "" || len(album.Artists) == 0 {
		jsonError(w, "Title and Artist are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 查找是否存在相同的专辑
	filter := bson.M{
		"title":   album.Title,
		"artists": album.Artists,
	}

	var existingAlbum Album
	err := albumCollection.FindOne(ctx, filter).Decode(&existingAlbum)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 不存在则插入新记录
			_, err = albumCollection.InsertOne(ctx, album)
			if err != nil {
				jsonError(w, "Failed to insert album", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(SuccessResponse{Message: "Album created successfully"})
			return
		}
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 构建更新文档，只更新非空字段
	update := bson.M{"$set": bson.M{}}
	if album.Genre != "" {
		update["$set"].(bson.M)["genre"] = album.Genre
	}
	update["$set"].(bson.M)["year"] = album.Year.Time
	// 在 handleMusicPost 函数中，修改 cuts 的处理逻辑：
	if album.Url != "" {
		update["$set"].(bson.M)["url"] = album.Url
	}
	if album.Artwork != "" {
		update["$set"].(bson.M)["artwork"] = album.Artwork
	}
	if album.Comment != "" {
		update["$set"].(bson.M)["comment"] = album.Comment
	}
	if album.Rating != 0 {
		update["$set"].(bson.M)["rating"] = album.Rating
	}

	// 如果有字段需要更新
	if len(update["$set"].(bson.M)) > 0 {
		_, err = albumCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			jsonError(w, "Failed to update album", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{Message: "Album updated successfully"})
}

func handleMusicDelete(w http.ResponseWriter, r *http.Request) {
	var album Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Music API endpoint called with album: %v", album)

	// 验证必需字段
	if album.Title == "" || len(album.Artists) == 0 {
		jsonError(w, "Title and Artist are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	filter := bson.M{
		"title":   album.Title,
		"artists": album.Artists,
	}

	result, err := albumCollection.DeleteOne(ctx, filter)
	if err != nil {
		jsonError(w, "Failed to delete album", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		jsonError(w, "Album not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{Message: "Album deleted successfully"})
}
