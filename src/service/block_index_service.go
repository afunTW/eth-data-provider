package service

import (
	"context"
	"math/big"

	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

type BlockIndexService struct {
	config          *config.Config
	client          *ethclient.Client
	blockNumberChan chan *big.Int
}

func NewBlockIndexService(config *config.Config) *BlockIndexService {
	client, err := ethclient.Dial(config.EthereumHost)
	if err != nil {
		log.Fatalf("BlockIndexService: failed %v\n", err)
	}
	return &BlockIndexService{
		config:          config,
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
			// TODO: get block, tx, receipt, log and store to db
			block, err := s.client.BlockByNumber(ctx, blockNum)
			if err != nil {
				log.Errorf("BlockWorker(blockNum=%v): failed %v\n", blockNum, err)
				continue
			}
			transactions := block.Transactions()
			log.Debugf("BlockWorker(blockNum=%v): get %v transactions\n", blockNum, len(transactions))
		}
	}
}
