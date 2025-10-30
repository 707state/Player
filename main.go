package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	albumCollectionName   = "albums"
	albumIndexName        = "album_title_artist_index"
	singlesCollectionName = "singles"
	singlesIndexName      = "album_artists_title_index"
	bookCollectionName    = "books"
	bookIndexName         = "book_title_author_index"
	filmCollectionName    = "films"
	filmIndexName         = "film_title_director_index"
)

type IndexField struct {
	Key   string
	Value any
}

var albumCollection *mongo.Collection
var singlesCollection *mongo.Collection
var bookCollection *mongo.Collection
var filmCollection *mongo.Collection

//go:embed dist
var staticFiles embed.FS

func main() {
	dbName := getEnv("DB_NAME", "musaic_dev")
	client := getConnection()
	defer client.Disconnect(context.TODO())

	// Ensure collections exist and create indexes
	ctx := context.Background()
	database := client.Database(dbName) // Using default database

	// Create albums collection if not exists and create index
	albumCollection = createCollectionIfNotExists(ctx, database, albumCollectionName)
	createIndexForCollection(ctx, albumCollection, albumIndexName, []IndexField{
		{Key: "artists", Value: 1},
		{Key: "title", Value: 1},
	})
	singlesCollection = createCollectionIfNotExists(ctx, database, singlesIndexName)
	createIndexForCollection(ctx, singlesCollection, singlesIndexName, []IndexField{
		{Key: "title", Value: 1},
		{Key: "artists", Value: 1},
		{Key: "album", Value: 1},
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
	httpListenAddress := getEnv("ADDRESS", "0.0.0.0")
	httpListenPort := getEnv("PORT", "8080")

	distFS, _ := fs.Sub(staticFiles, "dist")
	// 构建静态文件服务
	staticServer := http.FileServer(http.FS(distFS))

	http.Handle("/", staticServer)

	go http.HandleFunc("/music", corsMiddleware(handleMusic))
	go http.HandleFunc("/books", corsMiddleware(handleBooks))
	go http.HandleFunc("/movies", corsMiddleware(handleMovies))
	go http.HandleFunc("/single", corsMiddleware(handleSingle))
	bindAddress := fmt.Sprintf("%s:%s", httpListenAddress, httpListenPort)
	log.Printf("Server listening on %s", bindAddress)
	err := http.ListenAndServe(bindAddress, nil)
	if err != nil {
		log.Printf("Failed to ListenAndServe: %v\n", err.Error())
	}
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
	mongo_url := getEnv("MONGO_URL", "192.168.237.1")
	mongo_port := getEnvInt("MONGO_PORT", 7899)
	mongo_user := getEnv("MONGO_USER", "jask")
	mongo_password := getEnv("MONGO_PASSWORD", "theonlylove145")
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
