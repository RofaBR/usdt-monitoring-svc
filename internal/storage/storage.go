package storage

import (
	"context"
)

type TransferEvent struct {
	From   string
	To     string
	Amount string
	TxHash string
}

/*
	type TransferEventFilter struct {
		From         *string
		To           *string
		Counterparty *string
		Amount      *string
	}
*/
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
