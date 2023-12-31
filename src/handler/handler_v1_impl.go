package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type handlerV1Impl struct {
	config *config.Config
}

func NewHandlerV1Impl(config *config.Config) Handler {
	return &handlerV1Impl{config: config}
}

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

func (h *handlerV1Impl) GetTransactionByTxHash(ctx *gin.Context) {}
