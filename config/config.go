package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type ServerConfig struct {
	ListenAddr string        `env:"SERVER_ADDR,default=0.0.0.0:8080"`
	Difficulty byte          `env:"DIFFICULTY,default=5"`
	KeepAlive  time.Duration `env:"SERVER_KEEP_ALIVE,default=10s"`
}

type ClientConfig struct {
	ServerAddr   string `env:"SERVER_ADDR,default=0.0.0.0:8080"`
	RequestCount int    `env:"CLIENT_REQUEST_COUNT,default=10"`
}

func NewConfig[C any](ctx context.Context, config C) (*C, error) {
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
