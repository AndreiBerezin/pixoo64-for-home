package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AndreiBerezin/pixoo64/internal/collector"
	"github.com/AndreiBerezin/pixoo64/internal/state"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	log.Init()
	defer log.Sync()

	if os.Getenv("PIXOO_ADDRESS") == "" {
		log.Fatal("PIXOO_ADDRESS is empty,please check your environment variables")
	}

	log.Info("app started")

	collector := collector.New()
	collector.Start()

	state := state.New(collector)
	state.Start()

	waitShutdownSignal()
}

func waitShutdownSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Info("received shutdown signal: " + sig.String())
}
