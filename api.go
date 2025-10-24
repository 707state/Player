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

type Album struct {
	Title  string   `json:"title" bson:"title"`
	Artist string   `json:"artist" bson:"artist"`
	Genre  string   `json:"genre" bson:"genre"`
	Year   int      `json:"year" bson:"year"`
	Cuts   []string `json:"cuts" bson:"cuts"`
	Url    string   `json:"url" bson:"url"`
	// url of artwork!
	Artwork string `json:"artwork" bson:"artwork"`
	Comment string `json:"comment" bson:"comment"`
	// between 0-5
	Rating int `json:"rating" bson:"rating"`
}
type Book struct {
	Title  string `json:"title" bson:"title"`
	Author string `json:"author" bson:"author"`
	Genre  string `json:"genre" bson:"genre"`
	Year   int    `json:"year" bson:"year"`
	Url    string `json:"url" bson:"url"`
	// url of cover image!
	Cover   string `json:"cover" bson:"cover"`
	Comment string `json:"comment" bson:"comment"`
	// between 0-5
	Rating int `json:"rating" bson:"rating"`
}
type Movie struct {
	Title    string `json:"title" bson:"title"`
	Director string `json:"director" bson:"director"`
	Genre    string `json:"genre" bson:"genre"`
	Year     int    `json:"year" bson:"year"`
	Url      string `json:"url" bson:"url"`
	Comment  string `json:"comment" bson:"comment"`
	Rating   int    `json:"rating" bson:"rating"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func handleMusic(w http.ResponseWriter, r *http.Request) {
	log.Printf("Music API endpoint called")
	switch r.Method {
	case http.MethodGet:
		handleMusicGet(w, r)
	case http.MethodPost:
		handleMusicPost(w, r)
	case http.MethodDelete:
		handleMusicDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func handleMusicGet(w http.ResponseWriter, r *http.Request) {

	// parse query params from URL and build a MongoDB filter
	q := r.URL.Query()
	filter := bson.M{}

	if v := strings.TrimSpace(q.Get("title")); v != "" {
		filter["title"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("artist")); v != "" {
		filter["artist"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("genre")); v != "" {
		filter["genre"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("url")); v != "" {
		filter["url"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("artwork")); v != "" {
		filter["artwork"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("comment")); v != "" {
		filter["comment"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("year")); v != "" {
		if y, err := strconv.Atoi(v); err == nil {
			filter["year"] = y
		}
	}
	if v := strings.TrimSpace(q.Get("rating")); v != "" {
		if rInt, err := strconv.Atoi(v); err == nil {
			filter["rating"] = rInt
		}
	}
	if v := strings.TrimSpace(q.Get("cuts")); v != "" {
		parts := []string{}
		for _, p := range strings.Split(v, ",") {
			if t := strings.TrimSpace(p); t != "" {
				parts = append(parts, t)
			}
		}
		if len(parts) > 0 {
			// match documents that contain all provided cuts
			filter["cuts"] = bson.M{"$all": parts}
		}
	}
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
	if album.Title == "" || album.Artist == "" {
		jsonError(w, "Title and Artist are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 查找是否存在相同的专辑
	filter := bson.M{
		"title":  album.Title,
		"artist": album.Artist,
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
	if album.Year != 0 {
		update["$set"].(bson.M)["year"] = album.Year
	}
	if len(album.Cuts) > 0 {
		update["$set"].(bson.M)["cuts"] = album.Cuts
	}
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
	if album.Title == "" || album.Artist == "" {
		jsonError(w, "Title and Artist are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	filter := bson.M{
		"title":  album.Title,
		"artist": album.Artist,
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

func handleBooks(w http.ResponseWriter, r *http.Request) {
	log.Printf("Books API endpoint called")
	switch r.Method {
	case http.MethodGet:
		handleBooksGet(w, r)
	case http.MethodPost:
		handleBooksPost(w, r)
	case http.MethodDelete:
		handleBooksDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func handleBooksGet(w http.ResponseWriter, r *http.Request) {

	// parse query params from URL and build a MongoDB filter
	q := r.URL.Query()
	filter := bson.M{}
	if v := strings.TrimSpace(q.Get("title")); v != "" {
		filter["title"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("author")); v != "" {
		filter["author"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("genre")); v != "" {
		filter["genre"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("url")); v != "" {
		filter["url"] = bson.M{"$regex": v, "$options": "i"}
	}
	if v := strings.TrimSpace(q.Get("cover")); v != "" {
		filter["cover"] = bson.M{"$regex": v, "$options": "i"}
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
	log.Printf("Books API endpoint called with filter: %v", filter)
	ctx := r.Context()
	cursor, err := bookCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)
	var books []Book
	if err = cursor.All(ctx, &books); err != nil {
		http.Error(w, "Failed to decode books", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, "Failed to encode books", http.StatusInternalServerError)
		return
	}

}
func handleBooksPost(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Books API endpoint called with book: %v", book)

	// 验证必需字段
	if book.Title == "" || book.Author == "" {
		jsonError(w, "Title and Author are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 查找是否存在相同的书籍
	filter := bson.M{
		"title":  book.Title,
		"author": book.Author,
	}

	var existingBook Book
	err := bookCollection.FindOne(ctx, filter).Decode(&existingBook)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 不存在则插入新记录
			_, err = bookCollection.InsertOne(ctx, book)
			if err != nil {
				jsonError(w, "Failed to insert book", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(SuccessResponse{Message: "Book created successfully"})
			return
		}
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 构建更新文档，只更新非空字段
	update := bson.M{"$set": bson.M{}}
	if book.Genre != "" {
		update["$set"].(bson.M)["genre"] = book.Genre
	}
	if book.Year != 0 {
		update["$set"].(bson.M)["year"] = book.Year
	}
	if book.Url != "" {
		update["$set"].(bson.M)["url"] = book.Url
	}
	if book.Cover != "" {
		update["$set"].(bson.M)["cover"] = book.Cover
	}
	if book.Comment != "" {
		update["$set"].(bson.M)["comment"] = book.Comment
	}
	if book.Rating != 0 {
		update["$set"].(bson.M)["rating"] = book.Rating
	}

	// 如果有字段需要更新
	if len(update["$set"].(bson.M)) > 0 {
		result, err := bookCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			jsonError(w, "Failed to update book", http.StatusInternalServerError)
			return
		}
		if result.ModifiedCount == 0 {
			jsonError(w, "No changes made to book", http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{Message: "Book updated successfully"})
}

func handleBooksDelete(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Books API endpoint called with book: %v", book)

	// 验证必需字段
	if book.Title == "" || book.Author == "" {
		jsonError(w, "Title and Author are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	filter := bson.M{
		"title":  book.Title,
		"author": book.Author,
	}

	result, err := bookCollection.DeleteOne(ctx, filter)
	if err != nil {
		jsonError(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		jsonError(w, "Book not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{Message: "Book deleted successfully"})
}

func handleMovies(w http.ResponseWriter, r *http.Request) {
	log.Printf("Movies API endpoint called")
	switch r.Method {
	case http.MethodGet:
		handleMoviesGet(w, r)
	case http.MethodPost:
		handleMoviesPost(w, r)
	case http.MethodDelete:
		handleMoviesDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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
