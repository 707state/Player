package db

import (
	"log"
	"net/http"
)

func HandleMusic(w http.ResponseWriter, r *http.Request) {
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

func HandleSingle(w http.ResponseWriter, r *http.Request) {
	log.Printf("AlbumSingle API endpoint called")
	switch r.Method {
	case http.MethodGet:
		handleAlbumSingleGet(w, r)
	case http.MethodPost:
		handleAlbumSinglePost(w, r)
	case http.MethodDelete:
		handleAlbumSingleDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleBooks(w http.ResponseWriter, r *http.Request) {
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

func HandleMovies(w http.ResponseWriter, r *http.Request) {
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
