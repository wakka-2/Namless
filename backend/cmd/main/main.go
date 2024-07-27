/*
Package main is the entry point for this application.
*/
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wakka-2/Namless/backend/pkg/api"
	"github.com/wakka-2/Namless/backend/pkg/configs"
	"github.com/wakka-2/Namless/backend/pkg/repository"
	"github.com/wakka-2/Namless/backend/pkg/service"
)

const (
	readHeaderTimeout = 3 * time.Minute
)

// @title           Data storage API
// @version         1.0
// @description     This is the data storage API.
// @host      localhost:8080
// main starts the application.
func main() {
	configLocation := flag.String("config", "/etc/data/recon.json", "`configfile` for data service.")
	flag.Parse()

	cfg, err := configs.ReadConfigs(*configLocation)
	if err != nil {
		panic(err)
	}

	obfuscated, err := cfg.Obfuscate()
	if err != nil {
		panic("could not obfuscate configs")
	}

	log.Default().Printf("Starting with these configs: %s", obfuscated)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := repository.New(cfg.DSN, true)
	if err != nil {
		panic("could not build repository")
	}

	dataService := service.New(ctx, database)

	restAPI := api.New(dataService)

	go runServer(restAPI, cfg.ListenAddress)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}

func runServer(restAPI *api.RESTAPI, listenAddress string) {
	routes := restAPI.BuildMultiplexer()

	log.Default().Printf("Ready, accepting REST calls on %q...", listenAddress)

	server := &http.Server{
		Addr:              listenAddress,
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           routes,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Default().Printf("ListenAndServe error: %s", err)
	}
}
