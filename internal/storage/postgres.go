package storage

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type PostgresStorage struct {
	db *pgdb.DB
}

func NewPostgresStorage(db *pgdb.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) DB() *pgdb.DB {
	return s.db
}

func (s *PostgresStorage) SaveTransferEvent(ctx context.Context, event TransferEvent) error {
	log.Printf("Amount before formatting: %s", event.Amount)
	event.Amount = formatAmount(event.Amount, 6)
	log.Printf("Amount after formatting: %s", event.Amount)

	query := squirrel.Insert("transfers").
		Columns("from_address", "to_address", "amount", "transaction_hash", "block_number", "timestamp").
		Values(event.From, event.To, event.Amount, event.TransactionHash, event.BlockNumber, event.Timestamp)

	err := s.db.Exec(query)
	if err != nil {
		log.Printf("Error saving transfer event (From: %s, To: %s, TransactionHash: %s): %v", event.From, event.To, event.TransactionHash, err)
		return fmt.Errorf("SaveTransferEvent failed: %w", err)
	}

	return nil
}

func (s *PostgresStorage) QueryTransfers(ctx context.Context, query squirrel.Sqlizer) ([]TransferEvent, error) {
	var events []TransferEvent
	err := s.db.Select(&events, query)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, err
	}

	for _, event := range events {
		log.Printf("Retrieved event: %+v", event)
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
