package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/makasim/flowstatesrv/internal/app"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	cfg := app.Config{
		Driver: "memdriver",
		BadgerDriver: app.BadgerDriverConfig{
			Path: "badgerdb",
		},
	}
	if os.Getenv("FLOWSTATESRV_DRIVER") != "" {
		cfg.Driver = os.Getenv("FLOWSTATESRV_DRIVER")
	}
	if os.Getenv("FLOWSTATESRV_BADGERDRIVER_PATH") != "" {
		cfg.BadgerDriver.Path = os.Getenv("FLOWSTATESRV_BADGERDRIVER_PATH")
	}
	if os.Getenv("FLOWSTATESRV_BADGERDRIVER_IN_MEMORY") != "" {
		cfg.BadgerDriver.InMemory = os.Getenv("FLOWSTATESRV_BADGERDRIVER_IN_MEMORY") == `true`
	}
	if os.Getenv("FLOWSTATESRV_PGDRIVER_CONN_STRING") != "" {
		cfg.PostgresDriver.ConnString = os.Getenv("FLOWSTATESRV_PGDRIVER_CONN_STRING")
	}

	if err := app.New(cfg).Run(ctx); err != nil {
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}

func handleSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	<-signals
	log.Printf("INFO: got signal; canceling context")
	cancel()

	<-signals
	log.Printf("WARN: got second signal; force exiting")
	os.Exit(1)
}
