package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"musaic/db"
	"musaic/util"
	"net/http"
)

//go:embed dist
var staticFiles embed.FS

func main() {
	client := db.Init()
	defer client.Disconnect(context.TODO())
	// http api
	httpListenAddress := util.GetEnv("ADDRESS", "0.0.0.0")
	httpListenPort := util.GetEnv("PORT", "9999")
	useTLS := util.GetEnvBool("USE_TLS", true)
	certFile := util.GetEnv("CERT_FILE", "./localhost.pem")
	keyFile := util.GetEnv("KEY_FILE", "./localhost-key.pem")

	distFS, _ := fs.Sub(staticFiles, "dist")
	// 构建静态文件服务
	staticServer := http.FileServer(http.FS(distFS))

	http.Handle("/", staticServer)

	http.HandleFunc("/music", util.CorsMiddleware(db.HandleMusic))
	http.HandleFunc("/books", util.CorsMiddleware(db.HandleBooks))
	http.HandleFunc("/movies", util.CorsMiddleware(db.HandleMovies))
	http.HandleFunc("/single", util.CorsMiddleware(db.HandleSingle))
	// http.HandleFunc("/llm", util.CorsMiddleware())
	bindAddress := fmt.Sprintf("%s:%s", httpListenAddress, httpListenPort)
	if useTLS {
		log.Printf("HTTPS server listening on %s (cert=%s, key=%s)", bindAddress, certFile, keyFile)
		err := http.ListenAndServeTLS(bindAddress, certFile, keyFile, nil)
		if err != nil {
			log.Printf("Failed to ListenAndServeTLS: %v\n", err.Error())
		}
		return
	}
	log.Printf("HTTP server listening on %s", bindAddress)
	err := http.ListenAndServe(bindAddress, nil)
	if err != nil {
		log.Printf("Failed to ListenAndServe: %v\n", err.Error())
	}
}
