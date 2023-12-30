package handler

import (
	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/gin-gonic/gin"
)

type handlerV1Impl struct {
	config *config.Config
}

func NewHandlerV1Impl(config *config.Config) Handler {
	return &handlerV1Impl{config: config}
}

func (h *handlerV1Impl) GetLatestNBlocks(g *gin.Context)       {}
func (h *handlerV1Impl) GetBlockById(g *gin.Context)           {}
func (h *handlerV1Impl) GetTransactionByTxHash(g *gin.Context) {}
