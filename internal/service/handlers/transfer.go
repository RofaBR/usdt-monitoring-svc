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
	filter, err := NewTransferFilter(r)
	if err != nil {
		log.Println("Failed to decode query parameters:", err)
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	log.Println("Decoded filter:", filter)

	queryBuilder := squirrel.Select("from_address", "to_address", "amount", "transaction_hash", "block_number", "timestamp").
		From("transfers").
		PlaceholderFormat(squirrel.Dollar)

	if filter.From != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"from_address": *filter.From})
	}
	if filter.To != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"to_address": *filter.To})
	}
	if filter.Counterparty != nil {
		queryBuilder = queryBuilder.Where(squirrel.Or{
			squirrel.Eq{"from_address": *filter.Counterparty},
			squirrel.Eq{"to_address": *filter.Counterparty},
		})
	}
	if filter.Amount != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"amount": *filter.Amount})
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
