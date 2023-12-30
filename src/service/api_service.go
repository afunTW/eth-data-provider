package service

import (
	"context"
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

func NewApiService(config *config.Config, router *gin.Engine) (service *ApiService) {
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
	}
}

func (s *ApiService) Run(ctx context.Context) error {
	log.Info("Run server")
	if s.router == nil {
		log.Fatal("Missing router")
	}

	// run server
	// log the unexpected error before it raise
	err := s.server.ListenAndServe()
	defer func() {
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("Unexpected http server error: %v\n", err)
		}
	}()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}
