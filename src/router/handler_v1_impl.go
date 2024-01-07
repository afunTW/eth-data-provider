package router

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/afunTW/eth-data-provider/src/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type handlerV1Impl struct {
	config *config.Config
	repo   repository.EthereumIndexRepository
}

func NewHandlerV1Impl(config *config.Config, repo repository.EthereumIndexRepository) Handler {
	return &handlerV1Impl{config: config, repo: repo}
}

// @summary Get Latest blocks
// @schemes https
// @tags Block
// @accept json
// @produce json
// @param request query ReqGetLatestBlocks true "query param"
// @success 200 {object} RespGetLatestBlocks
// @failure 400 {object} RespStatus
// @failure 500 {object} RespStatus
// @router /api/v1/blocks [get]
func (h *handlerV1Impl) GetLatestBlocks(ctx *gin.Context) {
	var query ReqGetLatestBlocks
	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, RespStatus{
			Message:   "Fail to parse the query",
			Timestamp: uint64(time.Now().UnixNano()),
		})
		return
	}
	// get latest N blocks
	blockRecords, err := h.repo.GetLatestBlock(int(query.Limit))
	if err != nil {
		log.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RespStatus{
			Message:   "Fail to fetch blocks",
			Timestamp: uint64(time.Now().UnixNano()),
		})
	}
	// transform
	var blockResponse []BlockInfo
	for _, blockRecord := range blockRecords {
		blockResponse = append(blockResponse, BlockInfo{
			BlockNum:   blockRecord.BlockNum,
			BlockHash:  blockRecord.BlockHash,
			BlockTime:  blockRecord.BlockTime,
			ParentHash: blockRecord.ParentHash,
		})
	}
	// result
	ctx.JSON(http.StatusOK, RespGetLatestBlocks{
		Blocks: blockResponse,
		RespStatus: RespStatus{
			Message:   "Success",
			Timestamp: uint64(time.Now().UnixNano()),
		},
	})
}

// @summary Get block by given block id
// @schemes https
// @tags Block
// @accept json
// @produce json
// @param id path int true "block id"
// @success 200 {object} RespGetBlockDetail
// @failure 400 {object} RespStatus
// @failure 500 {object} RespStatus
// @router /api/v1/blocks/{id} [get]
func (h *handlerV1Impl) GetBlockById(ctx *gin.Context) {
	argBlockId := ctx.Param("id")
	blockId, err := strconv.Atoi(argBlockId)
	if err != nil {
		log.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, RespStatus{
			Message:   fmt.Sprintf("Failed to get the block id %v", argBlockId),
			Timestamp: uint64(time.Now().UnixNano()),
		})
		return
	}
	// get block info by given block id
	blockRecord, err := h.repo.GetBlock(uint64(blockId))
	if err != nil {
		log.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RespStatus{
			Message:   "Fail to fetch blocks",
			Timestamp: uint64(time.Now().UnixNano()),
		})
	}
	transactionRecords, err := h.repo.GetTransactions(uint64(blockId))
	if err != nil {
		log.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RespStatus{
			Message:   "Fail to fetch transactions",
			Timestamp: uint64(time.Now().UnixNano()),
		})
	}
	// transform
	log.Debugf("Get blockId %v with %v transactions\n", blockId, len(transactionRecords))
	blockResponse := BlockInfo{
		BlockNum:   blockRecord.BlockNum,
		BlockHash:  blockRecord.BlockHash,
		BlockTime:  blockRecord.BlockTime,
		ParentHash: blockRecord.ParentHash,
	}
	var transactionsResponse []string
	for _, transactionRecord := range transactionRecords {
		transactionsResponse = append(transactionsResponse, transactionRecord.TxHash)
	}
	// result
	ctx.JSON(http.StatusOK, RespGetBlockDetail{
		BlockInfo:    blockResponse,
		Transactions: transactionsResponse,
		RespStatus: RespStatus{
			Message:   "Success",
			Timestamp: uint64(time.Now().UnixNano()),
		},
	})
}

// @summary Get transactions by given tx hash
// @schemes https
// @tags Block
// @accept json
// @produce json
// @param txHash path string true "tx Hash"
// @success 200 {object} RespGetTransactionDetail
// @failure 400 {object} RespStatus
// @failure 500 {object} RespStatus
// @router /api/v1/transaction/{txHash} [get]
func (h *handlerV1Impl) GetTransactionByTxHash(ctx *gin.Context) {
	argTxHash := ctx.Param("txHash")
	// get transaction info by given tx hash
	transactionRecord, err := h.repo.GetTransaction(argTxHash)
	if err != nil {
		log.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RespStatus{
			Message:   "Fail to fetch transactions",
			Timestamp: uint64(time.Now().UnixNano()),
		})
	}
	logRecords, err := h.repo.GetLogs(argTxHash)
	if err != nil {
		log.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RespStatus{
			Message:   "Fail to fetch logs",
			Timestamp: uint64(time.Now().UnixNano()),
		})
	}
	// transform
	log.Debugf("Get txHash %v with %v logs\n", argTxHash, len(logRecords))
	transactionResponse := TransactionInfo{
		TxHash: transactionRecord.TxHash,
		From:   transactionRecord.FromAddress,
		To:     transactionRecord.ToAddress,
		Nonce:  transactionRecord.Nonce,
		Data:  hex.EncodeToString(transactionRecord.Data),
		Value: transactionRecord.Value,
	}
	var logsResponse []TransactionEventLogInfo
	for _, logRecord := range logRecords {
		logsResponse = append(logsResponse, TransactionEventLogInfo{
			Index: logRecord.LogIndex,
			Data: hex.EncodeToString(logRecord.Data),
		})
	}
	// result
	ctx.JSON(http.StatusOK, RespGetTransactionDetail{
		TransactionInfo: transactionResponse,
		Logs:            logsResponse,
		RespStatus: RespStatus{
			Message:   "Success",
			Timestamp: uint64(time.Now().UnixNano()),
		},
	})
}
