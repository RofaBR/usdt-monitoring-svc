package storage

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) SaveTransferEvent(ctx context.Context, event TransferEvent) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO transfers (from_address, to_address, amount, tx_hash) VALUES ($1, $2, $3, $4)",
		event.From, event.To, event.Amount, event.TxHash)
	if err != nil {
		log.Println("Failed to save transfer event:", err)
		return err
	}
	return nil
}

//func (s *PostgresStorage) GetTransferEvents(ctx context.Context)
