package handler

import (
	"net/http"
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
	ctx.JSON(http.StatusOK, RespGetLatestBlocks{
		Blocks: []BlockInfo{},
		RespStatus: RespStatus{
			Message:   "Success",
			Timestamp: uint64(time.Now().UnixNano()),
		},
	})
}

func (h *handlerV1Impl) GetBlockById(ctx *gin.Context)           {}
func (h *handlerV1Impl) GetTransactionByTxHash(ctx *gin.Context) {}
