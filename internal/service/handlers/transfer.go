package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
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

	queryBuilder := squirrel.Select("from_address", "to_address", "amount", "transaction_hash").
		From("transfers").
		PlaceholderFormat(squirrel.Dollar)

	for param, field := range queryParams {
		value := r.URL.Query().Get(param)
		if value != "" {
			if param == "counterparty" {
				queryBuilder = queryBuilder.Where(squirrel.Or{
					squirrel.Eq{"from_address": value},
					squirrel.Eq{"to_address": value},
				})
			} else {
				queryBuilder = queryBuilder.Where(squirrel.Eq{field: value})
			}
		}
	}

	transfers, err := h.Storage.QueryTransfers(r.Context(), queryBuilder)
	if err != nil {
		http.Error(w, "Failed to get transfer events", http.StatusInternalServerError)
		log.Println("Failed to get transfer events:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response, err := json.MarshalIndent(transfers, "", "  ")
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Println("Failed to encode response:", err)
		return
	}

	if _, err := w.Write(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Println("Failed to write response:", err)
		return
	}
}
