package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
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

type Date struct {
	time.Time
}

const dateLayout = "2006-01-02"

func (d *Date) UnmarshalJSON(b []byte) error {
	var dateStr string
	if err := json.Unmarshal(b, &dateStr); err != nil {
		return err
	}

	t, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, expected yyyy-mm-dd: %w", err)
	}

	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(dateLayout))
}
