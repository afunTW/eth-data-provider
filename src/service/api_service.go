package service

import (
	"net/http"

	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ApiService struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
}

func NewApiService(config *config.Config, router *gin.Engine) (service *ApiService, err error) {
	if config.ServerBindAddr == "" {
		log.Fatal("Missing config `SERVER_BIND_ADDR`")
	}
	server := &http.Server{
		Addr:    config.ServerBindAddr,
		Handler: router,
	}
	return &ApiService{
		config: config,
		router: router,
		server: server,
	}, nil
}
