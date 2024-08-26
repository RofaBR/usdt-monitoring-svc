package service

import (
	"context"
	"math/big"

	"github.com/RofaBR/usdt-monitoring-svc/internal/config"
	usdt "github.com/RofaBR/usdt-monitoring-svc/internal/ethereum"
	"github.com/RofaBR/usdt-monitoring-svc/internal/storage"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	USDTContractAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	BlockRange          = 100
)

func (s *service) connectToEthereum(cfg config.Config) *ethclient.Client {
	rpcURL := cfg.RPCURL() + cfg.InfuraProjectID()
	client, err := ethclient.Dial(rpcURL)
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

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		s.log.WithError(err).Error("Failed to get latest block header")
		return
	}
	currentBlock := header.Number
	startBlock := new(big.Int).Sub(currentBlock, big.NewInt(BlockRange))

	for startBlock.Cmp(currentBlock) < 0 {
		endBlock := new(big.Int).Add(startBlock, big.NewInt(BlockRange))
		if endBlock.Cmp(currentBlock) > 0 {
			endBlock = currentBlock
		}

		query := ethereum.FilterQuery{
			FromBlock: startBlock,
			ToBlock:   endBlock,
			Addresses: []common.Address{contractAddress},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			s.log.WithError(err).Errorf("Failed to filter logs for block range [%s, %s]", startBlock.String(), endBlock.String())
			break
		}

		for _, vLog := range logs {
			transferEvent, err := usdtContract.ParseTransfer(vLog)
			if err != nil {
				s.log.WithError(err).Error("Failed to parse log")
				continue
			}

			s.log.Infof("Transfer event: From: %s, To: %s, Tokens: %s", transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Value.String())

			err = s.storage.SaveTransferEvent(context.Background(), storage.TransferEvent{
				From:   transferEvent.From.Hex(),
				To:     transferEvent.To.Hex(),
				Amount: transferEvent.Value.String(),
				TxHash: vLog.TxHash.Hex(),
			})
			if err != nil {
				s.log.WithError(err).Error("Failed to save transfer event")
			}
		}

		startBlock = new(big.Int).Add(endBlock, big.NewInt(1))
	}
}
