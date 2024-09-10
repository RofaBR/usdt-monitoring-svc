package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func writeError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"errors": []map[string]interface{}{
			{
				"status": fmt.Sprintf("%d", statusCode),
				"title":  http.StatusText(statusCode),
				"detail": errMsg,
			},
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to encode error response:", err)
	}
}
