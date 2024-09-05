package storage

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type TransferEvent struct {
	From            string    `db:"from_address" json:"from"`
	To              string    `db:"to_address" json:"to"`
	Amount          string    `db:"amount" json:"amount"`
	TransactionHash string    `db:"transaction_hash" json:"transaction_hash"`
	BlockNumber     uint64    `db:"block_number" json:"block_number"`
	Timestamp       time.Time `db:"timestamp" json:"timestamp"`
}

type FilterCondition struct {
	Field string
	Value interface{}
}

type TransferEventFilter struct {
	Conditions []FilterCondition
}

type Storage interface {
	DB() *pgdb.DB
	SaveTransferEvent(ctx context.Context, event TransferEvent) error
	QueryTransfers(ctx context.Context, query squirrel.Sqlizer) ([]TransferEvent, error)
	GetLastProcessedBlock(ctx context.Context) (uint64, error)
}
