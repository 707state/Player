package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	albumCollectionName = "albums"
	albumIndexName      = "album_title_artist_index"
	bookCollectionName  = "books"
	bookIndexName       = "book_title_author_index"
	filmCollectionName  = "films"
	filmIndexName       = "film_title_director_index"
)

type IndexField struct {
	Key   string
	Value any
}

var albumCollection *mongo.Collection
var bookCollection *mongo.Collection
var filmCollection *mongo.Collection

func main() {
	dbName := getEnv("DB_NAME", "musaic")
	client := getConnection()
	defer client.Disconnect(context.TODO())

	// Ensure collections exist and create indexes
	ctx := context.Background()
	database := client.Database(dbName) // Using default database

	// Create albums collection if not exists and create index
	albumCollection = createCollectionIfNotExists(ctx, database, albumCollectionName)
	createIndexForCollection(ctx, albumCollection, albumIndexName, []IndexField{
		{Key: "artist", Value: 1},
		{Key: "title", Value: 1},
	})
	// Create books collection if not exists and create index
	bookCollection = createCollectionIfNotExists(ctx, database, bookCollectionName)
	createIndexForCollection(ctx, bookCollection, bookIndexName, []IndexField{
		{Key: "title", Value: 1},
		{Key: "author", Value: 1},
	})

	// Create films collection if not exists and create index
	filmCollection = createCollectionIfNotExists(ctx, database, filmCollectionName)
	createIndexForCollection(ctx, filmCollection, filmIndexName, []IndexField{
		{Key: "title", Value: 1},
		{Key: "director", Value: 1},
	})
	log.Println("Indexes created successfully")
	// http api
	httpListenAddress := getEnv("ADDRESS", "localhost")
	httpListenPort := getEnv("PORT", "8080")
	http.HandleFunc("/music", handleMusic)
	http.HandleFunc("/books", handleBooks)
	http.HandleFunc("/movies", handleMovies)
	bindAddress := fmt.Sprintf("%s:%s", httpListenAddress, httpListenPort)
	log.Printf("Server listening on %s", bindAddress)
	http.ListenAndServe(bindAddress, nil)
}

func createCollectionIfNotExists(ctx context.Context, db *mongo.Database, collectionName string) *mongo.Collection {
	// In MongoDB, collections are created automatically when first used.
	// But we can explicitly create them using CreateCollection
	err := db.CreateCollection(ctx, collectionName)
	if err != nil {
		// If collection already exists, this will throw an error which we can ignore
		log.Printf("Collection %s already exists or error creating: %v", collectionName, err)
	} else {
		log.Printf("Collection %s created successfully", collectionName)
	}
	return db.Collection(collectionName)
}
func createIndexForCollection(ctx context.Context, collection *mongo.Collection, indexName string, fields []IndexField) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field.Key, Value: field.Value})
	}
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true).SetName(indexName),
	})
	return err
}

func getConnection() *mongo.Client {
	mongo_url := getEnv("MONGO_URL", "localhost")
	mongo_port := getEnvInt("MONGO_PORT", 27017)
	mongo_user := getEnv("MONGO_USER", "admin")
	mongo_password := getEnv("MONGO_PASSWORD", "password")
	clientOptions := options.Client()
	clientOptions.ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s:%d", mongo_user, mongo_password, mongo_url, mongo_port),
	)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	return client
}
