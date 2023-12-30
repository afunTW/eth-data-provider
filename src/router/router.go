package router

import (
	"github.com/afunTW/eth-data-provider/src/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(v1Handler handler.Handler) *gin.Engine {
	router := gin.New()

	// set customized router config
	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:  []string{"*"},
				AllowMethods:  []string{"GET"},
				AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
				ExposeHeaders: []string{"Content-Length"},
			},
		),
	)

	// set route
	v1 := router.Group("/api/v1")
	{
		v1.GET("/blocks", v1Handler.GetLatestNBlocks)
		v1.GET("/blocks/:id", v1Handler.GetBlockById)
		v1.GET("/transaction/:txHash", v1Handler.GetTransactionByTxHash)
	}

	return router
}
