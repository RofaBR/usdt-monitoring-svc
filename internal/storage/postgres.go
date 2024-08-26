package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	_, err := s.db.ExecContext(ctx, "INSERT INTO transfers (from_address, to_address, amount, transaction_hash) VALUES ($1, $2, $3, $4)",
		event.From, event.To, event.Amount, event.TxHash)
	if err != nil {
		log.Printf("Error saving transfer event (From: %s, To: %s, TxHash: %s): %v", event.From, event.To, event.TxHash, err)
		return fmt.Errorf("SaveTransferEvent failed: %w", err)
	}
	return nil
}

func (s *PostgresStorage) GetTransferEvents(ctx context.Context, filter TransferEventFilter) ([]TransferEvent, error) {
	query := "SELECT from_address, to_address, amount, tx_hash FROM transfers WHERE 1=1"
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
		if err := rows.Scan(&event.From, &event.To, &event.Amount, &event.TxHash); err != nil {
			log.Println("Failed to scan row:", err)
			return nil, err
		}
		event.Amount = formatAmount(event.Amount)
		events = append(events, event)
	}

	return events, nil
}

func formatAmount(amount string) string {
	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Println("Failed to parse amount:", err)
		return amount
	}
	return fmt.Sprintf("%.2f", value/100)
}
