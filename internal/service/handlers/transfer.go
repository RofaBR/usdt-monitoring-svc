package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RofaBR/usdt-monitoring-svc/internal/storage"
)

type Handler struct {
	Storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{Storage: storage}
}

func (h *Handler) GetTransfers(w http.ResponseWriter, r *http.Request) {

	queryParams := map[string]string{
		"from":         "from_address",
		"to":           "to_address",
		"counterparty": "counterparty",
		"amount":       "amount",
	}

	var filter storage.TransferEventFilter
	for param, field := range queryParams {
		value := r.URL.Query().Get(param)
		if value != "" {
			if param == "counterparty" {
				filter.Conditions = append(filter.Conditions, storage.FilterCondition{Field: "from_address", Value: value})
				filter.Conditions = append(filter.Conditions, storage.FilterCondition{Field: "to_address", Value: value})
			} else {
				filter.Conditions = append(filter.Conditions, storage.FilterCondition{Field: field, Value: value})
			}
		}
	}

	transfers, err := h.Storage.GetTransferEvents(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to get transfer events", http.StatusInternalServerError)
		log.Println("Failed to get transfer events:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transfers); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Println("Failed to encode response:", err)
		return
	}
}
