package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/ed16/word-of-wisdom/internal/app/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	cfg, err := config.NewConfig(ctx, config.ServerConfig{})
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(cfg)
	srv.Start(ctx)
}
