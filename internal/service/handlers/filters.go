package handlers

import (
	"log"
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type TransferFilter struct {
	From         *string `filter:"filter[from]"`
	To           *string `filter:"filter[to]"`
	Counterparty *string `filter:"filter[counterparty]"`
	Amount       *string `filter:"filter[amount]"`
}

func NewTransferFilter(r *http.Request) (TransferFilter, error) {
	var filter TransferFilter
	log.Println("URL Query:", r.URL.Query())
	err := urlval.DecodeSilently(r.URL.Query(), &filter)
	if err != nil {
		log.Println("Failed to decode query parameters silently:", err)
		return filter, err
	}
	log.Println("Decoded filter silently:", filter)
	return filter, nil
}
