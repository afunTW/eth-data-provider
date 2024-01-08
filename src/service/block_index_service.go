package service

import (
	"context"
	"math/big"
	"time"

	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/afunTW/eth-data-provider/src/repository"
	"github.com/afunTW/eth-data-provider/src/worker"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

type BlockIndexService struct {
	config        *config.Config
	repo          repository.EthereumIndexRepository
	client        *ethclient.Client
	blockNumberCh chan *big.Int
	closeCh       chan int8
	workerMap     map[int]worker.BlockWorker
}

func NewBlockIndexService(config *config.Config, repoBlockIndex repository.EthereumIndexRepository) *BlockIndexService {
	client, err := ethclient.Dial(config.EthereumHost)
	if err != nil {
		log.Fatalf("BlockIndexService: failed %v\n", err)
	}
	return &BlockIndexService{
		config:        config,
		repo:          repoBlockIndex,
		client:        client,
		blockNumberCh: make(chan *big.Int, config.EthereumBlockWorkerCount),
		closeCh:       make(chan int8, 1),
		workerMap:     make(map[int]worker.BlockWorker),
	}
}

func (s *BlockIndexService) Start(ctx context.Context) {
	log.Info("BlockIndexService: start")
	// init the latest block number
	head, err := s.client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("BlockIndexService: start failed %v\n", err)
	}
	log.Debugf("BlockIndexService: get the latest block number %v\n", head.Number.String())

	// init N workers to crawl block infomation till the latest block
	for i := 0; i < s.config.EthereumBlockWorkerCount; i++ {
		blockWorker := worker.NewBlockWorker(i, s.client, s.repo, s.blockNumberCh)
		s.workerMap[i] = blockWorker
		go blockWorker.Start(ctx)
	}
	startBlockNum := big.NewInt(int64(s.config.EthereumBlockInitFrom))
	endBlockNum := big.NewInt(0).Sub(head.Number, big.NewInt(int64(s.config.EthereumBlockConfirmCount)))
	s.spawnBlockWorker(startBlockNum, endBlockNum)

	// keep monitor the latest N blocks
	ticker := time.NewTicker(time.Duration(s.config.BlockPullIntervalMillisecond) * time.Millisecond)
	defer ticker.Stop()
	latestBlockNum := new(big.Int).Set(endBlockNum)
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.closeCh:
			return
		case <-ticker.C:
			head, err := s.client.HeaderByNumber(ctx, nil)
			if err != nil {
				log.Errorf("BlockIndexService: start monitor failed %v\n", err)
				continue
			}
			log.Debugf("BlockIndexService: get the latest block number %v\n", head.Number.String())
			endBlockNum = big.NewInt(0).Sub(head.Number, big.NewInt(int64(s.config.EthereumBlockConfirmCount)))
			if latestBlockNum.Cmp(endBlockNum) < 0 {
				s.spawnBlockWorker(latestBlockNum, endBlockNum)
				latestBlockNum = new(big.Int).Set(endBlockNum)
			}
		}
	}
}

func (s *BlockIndexService) Stop() {
	log.Info("BlockIndexService: stop")
	s.closeCh <- int8(1)
	for _, w := range s.workerMap {
		w.Stop()
	}
}

func (s *BlockIndexService) spawnBlockWorker(startBlockNum *big.Int, endBlockNum *big.Int) {
	for n := startBlockNum; n.Cmp(endBlockNum) <= 0; n.Add(n, big.NewInt(1)) {
		s.blockNumberCh <- new(big.Int).Set(n)
	}
}
