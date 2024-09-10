package handlers

import (
	"log"
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type TransferFilter struct {
	From         *string `filter:"from"`
	To           *string `filter:"to"`
	Counterparty *string `filter:"counterparty"`
	Amount       *string `filter:"amount"`
}

func NewTransferFilter(r *http.Request) (TransferFilter, error) {
	var filter TransferFilter
	log.Println("URL Query:", r.URL.Query())
	err := urlval.Decode(r.URL.Query(), &filter)
	if err != nil {
		log.Println("Failed to decode query parameters silently:", err)
		return filter, err
	}
	log.Println("Decoded filter silently:", filter)

	log.Printf("Filter values - From: %v, To: %v, Counterparty: %v, Amount: %v\n",
		filter.From, filter.To, filter.Counterparty, filter.Amount)

	return filter, nil
}
