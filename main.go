package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/afunTW/eth-data-provider/src/config"
	"github.com/afunTW/eth-data-provider/src/repository"
	"github.com/afunTW/eth-data-provider/src/router"
	"github.com/afunTW/eth-data-provider/src/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Start service")

	// init dependency
	config := config.NewConfig()
	ll, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(ll)
	v1Handler := router.NewHandlerV1Impl(config)
	router := router.NewRouter(v1Handler)
	repoBlockIndex := repository.NewEthereumIndexGormImpl(config.GetDsn())

	// run service
	apiService := service.NewApiService(config, router)
	go apiService.Run(ctx)
	blockIndexService := service.NewBlockIndexService(config, repoBlockIndex)
	go blockIndexService.Start(ctx)

	// block until recv signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("End service")
}

func PrettyPrint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Error(err)
	}
	return string(b)
}
