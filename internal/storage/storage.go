package storage

import (
	"context"
	"time"
)

type TransferEvent struct {
	From            string
	To              string
	Amount          string
	TransactionHash string
	BlockNumber     uint64
	Timestamp       time.Time
}

type FilterCondition struct {
	Field string
	Value interface{}
}

type TransferEventFilter struct {
	Conditions []FilterCondition
}

type Storage interface {
	SaveTransferEvent(ctx context.Context, event TransferEvent) error
	GetTransferEvents(ctx context.Context, filter TransferEventFilter) ([]TransferEvent, error)
}
