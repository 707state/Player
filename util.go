package main

import (
	"net/http"
	"os"
	"strconv"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if value, err := strconv.Atoi(value); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value, err := strconv.ParseBool(value); err == nil {
			return value
		}
	}
	return defaultValue
}
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}
func removeStringInPlace(slice *[]string, target string) {
	for i := 0; i < len(*slice); i++ {
		if (*slice)[i] == target {
			// 删除元素并移动指针
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			i-- // 调整索引
		}
	}
}
