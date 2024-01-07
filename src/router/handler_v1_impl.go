package router

import (
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
	// TODO: get latest N blocks
	ctx.JSON(http.StatusOK, RespGetLatestBlocks{
		Blocks: []BlockInfo{},
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
	// TODO: get block info by given block id
	log.Debugf("Get blockId %v\n", blockId)
	ctx.JSON(http.StatusOK, RespGetBlockDetail{
		BlockInfo:    BlockInfo{},
		Transactions: []string{},
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
	// TODO: get transaction info by given tx hash
	log.Debugf("Get txHash %v\n", argTxHash)
	ctx.JSON(http.StatusOK, RespGetTransactionDetail{
		TransactionInfo: TransactionInfo{},
		Logs:            []TransactionEventLogInfo{},
		RespStatus: RespStatus{
			Message:   "Success",
			Timestamp: uint64(time.Now().UnixNano()),
		},
	})
}
