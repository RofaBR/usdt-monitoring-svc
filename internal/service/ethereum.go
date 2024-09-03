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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	transferEventSignature = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	BlockRange             = 10
)

func (s *service) connectToEthereum(cfg config.Config) *ethclient.Client {
	rpcURL := cfg.RPCURL()
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

func (s *service) getDecimalsFactor(usdtContract *usdt.Usdt) *big.Int {
	decimals, err := usdtContract.Decimals(nil)
	if err != nil {
		s.log.WithError(err).Error("Failed to get token decimals")
		return nil
	}
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
}

func (s *service) getCurrentBlock(client *ethclient.Client) (*big.Int, error) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		s.log.WithError(err).Error("Failed to get latest block header")
		return nil, err
	}
	return header.Number, nil
}

func (s *service) calculateEndBlock(startBlock, currentBlock *big.Int) *big.Int {
	endBlock := new(big.Int).Add(startBlock, big.NewInt(BlockRange))
	if endBlock.Cmp(currentBlock) > 0 {
		endBlock = currentBlock
	}
	return endBlock
}

func (s *service) filterLogs(client *ethclient.Client, startBlock, endBlock *big.Int, contractAddress common.Address) []types.Log {
	query := ethereum.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash(transferEventSignature)}},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to filter logs for block range [%s, %s]", startBlock.String(), endBlock.String())
		return nil
	}

	return logs
}

func (s *service) processLog(vLog types.Log, usdtContract *usdt.Usdt, decimalsFactor *big.Int) {
	transferEvent, err := usdtContract.ParseTransfer(vLog)
	if err != nil {
		s.log.WithError(err).Error("Failed to parse log")
		return
	}

	amount := new(big.Float).SetInt(transferEvent.Value)
	amount.Quo(amount, new(big.Float).SetInt(decimalsFactor))
	formattedAmount := amount.Text('f', 18)

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

	decimalsFactor := s.getDecimalsFactor(usdtContract)
	if decimalsFactor == nil {
		return
	}

	var startBlock *big.Int

	for {
		currentBlock, err := s.getCurrentBlock(client)
		if err != nil {
			return
		}

		if startBlock == nil {
			startBlock = new(big.Int).Sub(currentBlock, big.NewInt(BlockRange))
		}

		s.log.Infof("Starting to process blocks from %s to %s", startBlock.String(), currentBlock.String())

		for startBlock.Cmp(currentBlock) <= 0 {
			endBlock := s.calculateEndBlock(startBlock, currentBlock)

			s.log.Infof("Processing blocks from %s to %s", startBlock.String(), endBlock.String())

			logs := s.filterLogs(client, startBlock, endBlock, contractAddress)
			if logs == nil {
				s.log.Warn("No logs found or error occurred while filtering logs")
				break
			}

			s.log.Infof("Found %d logs in block range [%s, %s]", len(logs), startBlock.String(), endBlock.String())

			for _, vLog := range logs {
				s.log.Infof("Processing log with signature: %s", vLog.Topics[0].Hex())
				if vLog.Topics[0].Hex() != transferEventSignature {
					s.log.Infof("Skipping non-Transfer event with signature: %s", vLog.Topics[0].Hex())
					continue
				}

				s.processLog(vLog, usdtContract, decimalsFactor)
			}

			s.log.Infof("Finished processing block range [%s, %s]", startBlock.String(), endBlock.String())

			startBlock = new(big.Int).Add(endBlock, big.NewInt(1))

			time.Sleep(1 * time.Second)
		}

		s.log.Info("Completed processing all blocks. Waiting for new blocks...")

		startBlock = new(big.Int).Add(currentBlock, big.NewInt(1))

		time.Sleep(10 * time.Second)
	}
}
