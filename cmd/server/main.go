package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/ed16/word-of-wisdom/internal/app/server"
	"github.com/ed16/word-of-wisdom/pkg/tcp"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	cfg, err := config.NewConfig(ctx, config.ServerConfig{})
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(cfg, &tcp.DefaultConnector{})
	srv.Start(ctx)
}
