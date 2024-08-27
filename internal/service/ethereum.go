package service

import (
	"context"
	"math/big"
	"time"

	"github.com/RofaBR/usdt-monitoring-svc/internal/config"
	usdt "github.com/RofaBR/usdt-monitoring-svc/internal/ethereum"
	"github.com/RofaBR/usdt-monitoring-svc/internal/storage"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	transferEventSignature = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	BlockRange             = 100
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

	decimals, err := usdtContract.Decimals(nil)
	if err != nil {
		s.log.WithError(err).Error("Failed to get token decimals")
		return
	}
	decimalsFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

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

			if vLog.Topics[0].Hex() != transferEventSignature {
				s.log.Infof("Skipping non-Transfer event with signature: %s", vLog.Topics[0].Hex())
				continue
			}

			transferEvent, err := usdtContract.ParseTransfer(vLog)
			if err != nil {
				s.log.WithError(err).Error("Failed to parse log")
				continue
			}

			amount := new(big.Rat).SetInt(transferEvent.Value)
			amount.Quo(amount, new(big.Rat).SetInt(decimalsFactor))
			formattedAmount := amount.FloatString(int(decimals))

			s.log.Infof("Transfer event: From: %s, To: %s, Amount: %s, Value: %s", transferEvent.From.Hex(), transferEvent.To.Hex(), formattedAmount, transferEvent.Value.String())

			err = s.storage.SaveTransferEvent(context.Background(), storage.TransferEvent{
				From:            transferEvent.From.Hex(),
				To:              transferEvent.To.Hex(),
				Amount:          formattedAmount,
				TransactionHash: vLog.TxHash.Hex(),
				BlockNumber:     vLog.BlockNumber,
				Timestamp:       time.Now().UTC(),
			})
			if err != nil {
				s.log.WithError(err).Error("Failed to save transfer event")
			}
		}

		startBlock = new(big.Int).Add(endBlock, big.NewInt(1))
	}
}
