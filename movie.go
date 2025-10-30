package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Movie struct {
	Title    string `json:"title" bson:"title"`
	Director string `json:"director" bson:"director"`
	Genre    string `json:"genre" bson:"genre"`
	Year     int    `json:"year" bson:"year"`
	Url      string `json:"url" bson:"url"`
	Comment  string `json:"comment" bson:"comment"`
	Rating   int    `json:"rating" bson:"rating"`
}

func handleMoviesGet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := bson.M{}
	if v := strings.TrimSpace(q.Get("title")); v != "" {
		filter["title"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("director")); v != "" {
		filter["director"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("genre")); v != "" {
		filter["genre"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("url")); v != "" {
		filter["url"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("comment")); v != "" {
		filter["comment"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("year")); v != "" {
		y, err := strconv.Atoi(v)
		if err == nil {
			filter["year"] = y
		}
	}
	if v := strings.TrimSpace(q.Get("rating")); v != "" {
		if rInt, err := strconv.Atoi(v); err == nil {
			filter["rating"] = rInt
		}
	}
	log.Printf("Movies API endpoint called with filter: %v", filter)
	ctx := r.Context()
	cursor, err := filmCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)
	var movies []Movie
	if err = cursor.All(ctx, &movies); err != nil {
		http.Error(w, "Failed to decode movies", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		http.Error(w, "Failed to encode movies", http.StatusInternalServerError)
		return
	}
}

func handleMoviesPost(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Movies API endpoint called with movie: %v", movie)
	if movie.Title == "" || movie.Director == "" {
		jsonError(w, "Title and Director are required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	filter := bson.M{
		"title":    movie.Title,
		"director": movie.Director,
	}
	var existingMovie Movie
	err := filmCollection.FindOne(ctx, filter).Decode(&existingMovie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = filmCollection.InsertOne(ctx, movie)
			if err != nil {
				jsonError(w, "Failed to insert movie", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(SuccessResponse{Message: "Movie created successfully"})
			return
		}
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	update := bson.M{"$set": bson.M{}}
	if movie.Genre != "" {
		update["$set"].(bson.M)["genre"] = movie.Genre
	}
	if movie.Year != 0 {
		update["$set"].(bson.M)["year"] = movie.Year
	}
	if movie.Url != "" {
		update["$set"].(bson.M)["url"] = movie.Url
	}
	if movie.Comment != "" {
		update["$set"].(bson.M)["comment"] = movie.Comment
	}
	if movie.Rating != 0 {
		update["$set"].(bson.M)["rating"] = movie.Rating
	}
	if len(update["$set"].(bson.M)) > 0 {
		result, err := filmCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			jsonError(w, "Failed to update movie", http.StatusInternalServerError)
			return
		}
		if result.ModifiedCount == 0 {
			jsonError(w, "No changes made to movie", http.StatusNotModified)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{Message: "Movie updated successfully"})
}

func handleMoviesDelete(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Movies API endpoint called with movie: %v", movie)
	if movie.Title == "" || movie.Director == "" {
		jsonError(w, "Title and Director are required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	filter := bson.M{
		"title":    movie.Title,
		"director": movie.Director,
	}
	result, err := filmCollection.DeleteOne(ctx, filter)
	if err != nil {
		jsonError(w, "Failed to delete movie", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		jsonError(w, "Movie not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{Message: "Movie deleted successfully"})
}
