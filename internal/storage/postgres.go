package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) SaveTransferEvent(ctx context.Context, event TransferEvent) error {
	event.Amount = formatAmount(event.Amount, 6)
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO transfers (from_address, to_address, amount, transaction_hash, block_number, timestamp) VALUES ($1, $2, $3, $4, $5, $6)",
		event.From, event.To, event.Amount, event.TransactionHash, event.BlockNumber, event.Timestamp)
	if err != nil {
		log.Printf("Error saving transfer event (From: %s, To: %s, TransactionHash: %s): %v", event.From, event.To, event.TransactionHash, err)
		return fmt.Errorf("SaveTransferEvent failed: %w", err)
	}
	return nil
}

func (s *PostgresStorage) GetTransferEvents(ctx context.Context, filter TransferEventFilter) ([]TransferEvent, error) {
	query := "SELECT from_address, to_address, amount, transaction_hash FROM transfers WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	for _, condition := range filter.Conditions {
		query += fmt.Sprintf(" AND %s = $%d", condition.Field, argIndex)
		args = append(args, condition.Value)
		argIndex++
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println("Failed to get transfer events:", err)
		return nil, err
	}
	defer rows.Close()

	var events []TransferEvent
	for rows.Next() {
		var event TransferEvent
		if err := rows.Scan(&event.From, &event.To, &event.Amount, &event.TransactionHash); err != nil {
			log.Println("Failed to scan row:", err)
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func formatAmount(amount string, decimals int) string {
	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Println("Failed to parse amount:", err)
		return amount
	}
	factor := math.Pow(10, float64(decimals))
	return fmt.Sprintf("%.6f", value/factor)
}
