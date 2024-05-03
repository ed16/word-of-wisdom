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
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	cfg, err := config.NewConfig(ctx, config.ClientConfig{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	cliClient := client.NewClient(cfg)
	log.Println("Connecting to server...")

	for i := 0; i < cfg.RequestCount; i++ {
		quote, err := cliClient.RequestQuote(ctx)
		if err != nil {
			log.Printf("Failed to get quote: %s\n", err)
			return
		}
		log.Printf("Received quote #%d: %s\n", i+1, quote)
	}
}
