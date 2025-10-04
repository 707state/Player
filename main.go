package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func parseFiles(path string) <-chan map[string]any {
	ch := make(chan map[string]any)
	go func() {
		defer close(ch)

		cmd := exec.Command("exiftool",
			"-if", "$Filename !~ /^._/ and $FileType !~ /JPEG|PNG|GIF|BMP|TIFF/i",
			"-q", "-json", "-r", path,
		)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatalf("获取 stdout 出错: %v", err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatalf("启动 exiftool 出错: %v", err)
		}

		dec := json.NewDecoder(stdout)
		token, err := dec.Token()
		if err != nil {
			log.Fatalf("解析 JSON token 出错: %v", err)
		}
		if delim, ok := token.(json.Delim); !ok || delim != '[' {
			log.Fatalf("输出不是 JSON 数组")
		}

		for dec.More() {
			var fileInfo map[string]any
			if err := dec.Decode(&fileInfo); err != nil {
				log.Fatalf("解析 JSON 对象出错: %v", err)
			}
			ch <- fileInfo
		}

		if err := cmd.Wait(); err != nil {
			log.Fatalf("等待 exiftool 完成出错: %v", err)
		}
	}()
	return ch
}

func insertToMongo(ctx context.Context, uri string, ch <-chan map[string]any) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("MongoDB 连接失败: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("MongoDB 断开连接失败: %v", err)
		}
	}()

	// 获取数据库和 collection
	db := client.Database("musicdb")
	collection := db.Collection("tracks")
	for fileInfo := range ch {
		artist, _ := fileInfo["Artist"].(string)
		title, _ := fileInfo["Title"].(string)
		file_name, _ := fileInfo["FileName"].(string)
		// 如果 Artist 为空，就尝试从 Title 中推断
		if artist == "" && title != "" {
			parts := strings.SplitN(file_name, "-", 2)
			if len(parts) == 2 {
				artist = strings.TrimSpace(parts[0])
			}
		}
		doc := fileInfo
		if _, err := collection.InsertOne(ctx, doc); err != nil {
			log.Printf("插入 MongoDB 失败: %v", err)
		} else {
			fmt.Printf("已插入: %v - %v\n", artist, doc["Title"])
		}
	}
}

func main() {
	path := flag.String("path", "/Volumes/NO NAME/", "File path to music directory.")
	flag.Parse()
	fmt.Println("扫描路径: ", *path)

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		// uri = "mongodb://192.168.30.92:27017"
		log.Fatal("必须设置 'MONGODB_URI' 环境变量")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	ch := parseFiles(*path)
	insertToMongo(ctx, uri, ch)
}
