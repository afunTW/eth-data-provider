package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	GetLatestNBlocks(*gin.Context)
	GetBlockById(*gin.Context)
	GetTransactionByTxHash(*gin.Context)
}
