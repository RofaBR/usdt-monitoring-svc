package service

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func (s *service) connectToEthereum() *ethclient.Client {

	infuraProjectID := os.Getenv("INFURA_PROJECT_ID")

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + infuraProjectID)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	s.log.Info("Successfully connected to Ethereum")
	return client
}
