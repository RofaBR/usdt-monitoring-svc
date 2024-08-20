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

type Storage interface {
	SaveTransferEvent(ctx context.Context, event TransferEvent) error
	//GetTransferEvents(ctx context.Context)
}
