package service

import (
	"context"
	"math/big"

	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/afunTW/eth-data-provider/src/repository"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

type BlockIndexService struct {
	config          *config.Config
	repo            repository.EthereumIndexRepository
	client          *ethclient.Client
	blockNumberChan chan *big.Int
}

func NewBlockIndexService(config *config.Config, repoBlockIndex repository.EthereumIndexRepository) *BlockIndexService {
	client, err := ethclient.Dial(config.EthereumHost)
	if err != nil {
		log.Fatalf("BlockIndexService: failed %v\n", err)
	}
	return &BlockIndexService{
		config:          config,
		repo:            repoBlockIndex,
		client:          client,
		blockNumberChan: make(chan *big.Int, config.EthereumBlockWorkerCount),
	}
}

func (s *BlockIndexService) Start(ctx context.Context) {
	log.Info("BlockIndexService: start")
	// init the latest block number
	head, err := s.client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("BlockIndexService: start failed %v\n", err)
	}
	log.Debugf("BlockIndexService: get the latest block number %v\n", *head.Number)

	// init N workers to crawl block infomation till the latest block
	for i := 0; i < s.config.EthereumBlockWorkerCount; i++ {
		go s.blockWorker(ctx)
	}
	log.Debugf("BlockIndexService: spawn %v workers\n", s.config.EthereumBlockWorkerCount)
	startBlockNumber := big.NewInt(int64(s.config.EthereumBlockInitFrom))
	one := big.NewInt(1)
	for blockNum := startBlockNumber; blockNum.Cmp(head.Number) <= 0; blockNum.Add(blockNum, one) {
		s.blockNumberChan <- new(big.Int).Set(blockNum)
	}

	// TODO: keep monitor the latest N blocks
}

func (s *BlockIndexService) blockWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case blockNum := <-s.blockNumberChan:
			// handle block
			block, err := s.client.BlockByNumber(ctx, blockNum)
			if err != nil {
				log.Errorf("BlockWorker(blockNum=%v): failed %v\n", blockNum, err)
				continue
			}
			err = s.repo.AddBlocks([]*repository.EthereumBlock{
				{
					BlockNum:   block.Number().Uint64(),
					BlockHash:  block.Hash().Hex(),
					BlockTime:  block.Time(),
					ParentHash: block.ParentHash().Hex(),
				},
			})
			if err != nil {
				log.Errorf("BlockWorker(blockNum=%v): failed %v\n", blockNum, err)
			}
			// handle transactions
			// handle transactionLog from receipt
			txs := block.Transactions()
			log.Debugf("BlockWorker(blockNum=%v): get %v transactions\n", blockNum, len(txs))
			var txRecords []*repository.EthereumTransaction
			var logRecords []*repository.EthereumLog
			for _, tx := range txs {
				from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
				if err != nil {
					from, _ = types.Sender(types.HomesteadSigner{}, tx)
				}
				// TODO: ignore data if it excceed the limit of database blob
				txData := tx.Data()
				if len(txData) > 65535 {
					txData = nil
				}
				txRecords = append(txRecords, &repository.EthereumTransaction{
					TxHash:      tx.Hash().Hex(),
					BlockNum:    block.Number().Uint64(),
					FromAddress: from.Hex(),
					ToAddress:   tx.To().Hex(),
					Nonce:       tx.Nonce(),
					Data:        txData,
					Value:       tx.Value().Uint64(),
				})
				receipt, err := s.client.TransactionReceipt(ctx, tx.Hash())
				if err != nil {
					log.Errorf("BlockWorker(blockNum=%v): failed get receipt\n%v\n", blockNum, err)
				}
				for _, txLog := range receipt.Logs {
					if receipt.TxHash.Hex() != tx.Hash().Hex() {
						log.Warnf("receipt hash: %v\ntx hash: %v\n", receipt.TxHash.Hex(), tx.Hash().Hex())
					}
					txLogData := txLog.Data
					if len(txLog.Data) > 65535 {
						txLogData = nil
					}
					logRecords = append(logRecords, &repository.EthereumLog{
						TxHash:   receipt.TxHash.Hex(),
						LogIndex: uint64(txLog.Index),
						Data:     txLogData,
					})
				}
			}
			err = s.repo.AddTransactions(txRecords)
			if err != nil {
				log.Errorf("BlockWorker(blockNum=%v): failed add txs\n%v\n", blockNum, err)
			}
			log.Debugf("BlockWorker(blockNum=%v): processed %v transactions\n", blockNum, len(txRecords))
			err = s.repo.AddLogs(logRecords)
			if err != nil {
				log.Errorf("BlockWorker(blockNum=%v): failed add logs\n%v\n", blockNum, err)
			}
			log.Debugf("BlockWorker(blockNum=%v): processed %v logs\n", blockNum, len(logRecords))
		}
	}
}
