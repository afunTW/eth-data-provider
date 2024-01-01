package router

import "github.com/gin-gonic/gin"

type Handler interface {
	GetLatestBlocks(*gin.Context)
	GetBlockById(*gin.Context)
	GetTransactionByTxHash(*gin.Context)
}
