package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/afunTW/eth-data-provider/src/handler"
	"github.com/afunTW/eth-data-provider/src/router"
	"github.com/afunTW/eth-data-provider/src/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Start service")
	// run service in coroutine
	config := config.NewConfig()
	v1Handler := handler.NewHandlerV1Impl(config)
	router := router.NewRouter(v1Handler)
	service := service.NewApiService(config, router)
	go service.Run(ctx)

	// block until recv signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("End service")
}
