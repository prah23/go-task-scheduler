package views

import (
	"encoding/json"
	"net/http"
)

func Success(data interface{}, statusCode int, w http.ResponseWriter) {
	response := map[string]interface{}{
		"status": "success",
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func Error(data interface{}, statusCode int, w http.ResponseWriter) {
	response := map[string]interface{}{
		"status": "error",
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
