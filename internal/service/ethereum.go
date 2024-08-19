package service

import (
	"context"

	"github.com/RofaBR/usdt-monitoring-svc/internal/config"
	usdt "github.com/RofaBR/usdt-monitoring-svc/internal/ethereum"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const USDTContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"

func (s *service) connectToEthereum(cfg config.Config) *ethclient.Client {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + cfg.InfuraProjectID())
	if err != nil {
		s.log.WithError(err).Error("Failed to connect to the Ethereum client")
		return nil
	}

	s.log.Info("Successfully connected to Ethereum")
	return client
}

func (s *service) getContractAddress(cfg config.Config) common.Address {
	return common.HexToAddress(cfg.ContractAddress())
}

func (s *service) GetTransferEvents(cfg config.Config) {
	client := s.connectToEthereum(cfg)
	if client == nil {
		return
	}

	contractAddress := s.getContractAddress(cfg)
	usdtContract, err := usdt.NewUsdt(contractAddress, client)
	if err != nil {
		s.log.WithError(err).Error("Failed to instantiate a Token contract")
		return
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		s.log.WithError(err).Error("Failed to filter logs")
		return
	}

	for _, vLog := range logs {
		transferEvent, err := usdtContract.ParseTransfer(vLog)
		if err != nil {
			s.log.WithError(err).Error("Failed to parse log")
			continue
		}
		s.log.Infof("Transfer event: From: %s, To: %s, Tokens: %s", transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Value.String())
	}
}
