package db

import (
	"encoding/json"
	"log"
	"musaic/util"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, "Failed to encode books", http.StatusInternalServerError)
		return
	}

}
func handleBooksPost(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		util.JsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Books API endpoint called with book: %v", book)

	// 验证必需字段
	if book.Title == "" || book.Author == "" {
		util.JsonError(w, "Title and Author are required", http.StatusBadRequest)
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
				util.JsonError(w, "Failed to insert book", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(util.SuccessResponse{Message: "Book created successfully"})
			return
		}
		util.JsonError(w, "Database error", http.StatusInternalServerError)
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
			util.JsonError(w, "Failed to update book", http.StatusInternalServerError)
			return
		}
		if result.ModifiedCount == 0 {
			util.JsonError(w, "No changes made to book", http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.SuccessResponse{Message: "Book updated successfully"})
}

func handleBooksDelete(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		util.JsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Books API endpoint called with book: %v", book)

	// 验证必需字段
	if book.Title == "" || book.Author == "" {
		util.JsonError(w, "Title and Author are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	filter := bson.M{
		"title":  book.Title,
		"author": book.Author,
	}

	result, err := bookCollection.DeleteOne(ctx, filter)
	if err != nil {
		util.JsonError(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		util.JsonError(w, "Book not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(util.SuccessResponse{Message: "Book deleted successfully"})
}
