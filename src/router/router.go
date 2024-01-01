package router

import (
	"github.com/afunTW/eth-data-provider/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func NewRouter(v1Handler Handler) *gin.Engine {
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
		gin.Recovery(),
	)

	// set swagger info
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	docs.SwaggerInfo.Title = "Ethereum Data Service"

	// set root route
	root := router.Group("/")
	{
		root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// set v1 route
	v1 := router.Group("/api/v1")
	{
		v1.GET("/blocks", v1Handler.GetLatestBlocks)
		v1.GET("/blocks/:id", v1Handler.GetBlockById)
		v1.GET("/transaction/:txHash", v1Handler.GetTransactionByTxHash)
	}

	return router
}
