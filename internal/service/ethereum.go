package service

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const USDTContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"

func (s *service) connectToEthereum() *ethclient.Client {

	infuraProjectID := os.Getenv("INFURA_PROJECT_ID")

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + infuraProjectID)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	s.log.Info("Successfully connected to Ethereum")
	return client
}

func (s *service) loadABI() abi.ABI {
	abiFile := "internal/ethereum/usdt_abi.json"
	fileContent, err := os.ReadFile(abiFile)
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}

	var contractABI abi.ABI
	err = json.Unmarshal(fileContent, &contractABI)
	if err != nil {
		log.Fatalf("Failed to unmarshal ABI: %v", err)
	}

	return contractABI
}

func (s *service) getContractAddress() common.Address {
	return common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
}

func (s *service) GetTransferEvents(client *ethclient.Client) {
	contractAddress := s.getContractAddress()
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to filter logs: %v", err)
	}

	contractABI := s.loadABI()
	for _, vLog := range logs {
		transferEvent := struct {
			From   common.Address
			To     common.Address
			Tokens *big.Int
		}{}
		err := contractABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack log: %v", err)
		}
		log.Printf("Transfer event: From: %s, To: %s, Tokens: %s", transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Tokens.String())
	}
}
