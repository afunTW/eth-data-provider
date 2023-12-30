package router

import (
	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(config *config.Config) *gin.Engine {
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

	// TODO: set route

	return router
}
