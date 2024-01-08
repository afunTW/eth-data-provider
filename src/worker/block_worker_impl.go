package worker

import (
	"context"
	"math/big"

	"github.com/afunTW/eth-data-provider/src/repository"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

type blockWorkerImpl struct {
	id            int
	client        *ethclient.Client
	repo          repository.EthereumIndexRepository
	closeCh       chan int8
	blockNumberCh chan *big.Int
}

func NewBlockWorker(
	id int,
	client *ethclient.Client,
	repo repository.EthereumIndexRepository,
	blockNumberCh chan *big.Int,
) BlockWorker {
	return &blockWorkerImpl{
		id:            id,
		client:        client,
		repo:          repo,
		blockNumberCh: blockNumberCh,
		closeCh:       make(chan int8, 1),
	}
}

func (w *blockWorkerImpl) Start(ctx context.Context) {
	log.Infof("BlockWorker(id=%v): start\n", w.id)
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("BlockWorker(id=%v): panic %v\n", w.id, err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.closeCh:
			return
		case blockNum := <-w.blockNumberCh:
			// handle block
			block, err := w.client.BlockByNumber(ctx, blockNum)
			if err != nil {
				log.Errorf("BlockWorker(id=%v, blockNum=%v): failed %v\n", w.id, blockNum, err)
				continue
			}
			// TODO: check block nil pointer dereference problem
			err = w.repo.AddBlocks([]*repository.EthereumBlock{
				{
					BlockNum:   block.Number().Uint64(),
					BlockHash:  block.Hash().Hex(),
					BlockTime:  block.Time(),
					ParentHash: block.ParentHash().Hex(),
				},
			})
			if err != nil {
				log.Errorf("BlockWorker(id=%v, blockNum=%v): failed %v\n", w.id, blockNum, err)
			}
			// handle transactions
			// handle transactionLog from receipt
			txs := block.Transactions()
			log.Infof("BlockWorker(id=%v, blockNum=%v): get %v transactions\n", w.id, blockNum, len(txs))
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
				// TODO: check block nil pointer dereference problem
				txRecords = append(txRecords, &repository.EthereumTransaction{
					TxHash:      tx.Hash().Hex(),
					BlockNum:    block.Number().Uint64(),
					FromAddress: from.Hex(),
					ToAddress:   tx.To().Hex(),
					Nonce:       tx.Nonce(),
					Data:        txData,
					Value:       tx.Value().Uint64(),
				})
				receipt, err := w.client.TransactionReceipt(ctx, tx.Hash())
				if err != nil {
					log.Errorf("BlockWorker(id=%v, blockNum=%v): failed get receipt\n%v\n", w.id, blockNum, err)
				}
				for _, txLog := range receipt.Logs {
					if receipt.TxHash.Hex() != tx.Hash().Hex() {
						log.Warnf("receipt hash: %v\ntx hash: %v\n", receipt.TxHash.Hex(), tx.Hash().Hex())
					}
					txLogData := txLog.Data
					if len(txLog.Data) > 65535 {
						txLogData = nil
					}
					// TODO: check block nil pointer dereference problem
					logRecords = append(logRecords, &repository.EthereumLog{
						TxHash:   receipt.TxHash.Hex(),
						LogIndex: uint64(txLog.Index),
						Data:     txLogData,
					})
				}
			}
			err = w.repo.AddTransactions(txRecords)
			if err != nil {
				log.Errorf("BlockWorker(id=%v, blockNum=%v): failed add txs\n%v\n", w.id, blockNum, err)
			}
			log.Infof("BlockWorker(id=%v, blockNum=%v): processed %v transactions\n", w.id, blockNum, len(txRecords))
			err = w.repo.AddLogs(logRecords)
			if err != nil {
				log.Errorf("BlockWorker(id=%v, blockNum=%v): failed add logs\n%v\n", w.id, blockNum, err)
			}
			log.Infof("BlockWorker(id=%v, blockNum=%v): processed %v logs\n", w.id, blockNum, len(logRecords))
		}
	}
}

func (w *blockWorkerImpl) Stop() {
	log.Infof("BlockWorker(id=%v): stop", w.id)
	w.closeCh <- int8(1)
}
