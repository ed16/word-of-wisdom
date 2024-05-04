package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/ed16/word-of-wisdom/internal/app/client"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	go func() {
		<-ctx.Done()
		log.Println("Gracefully stopping...")
	}()

	cfg, err := config.NewConfig(ctx, config.ClientConfig{})
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	app := client.NewClient(cfg)
	log.Println("Connecting to server...")

	app.Start(ctx)

}
