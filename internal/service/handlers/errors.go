package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]string{
		"error": errMsg,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to encode error response:", err)
	}
}
