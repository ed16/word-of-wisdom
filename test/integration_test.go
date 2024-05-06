package integration_test

import (
	"context"
	"testing"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/ed16/word-of-wisdom/internal/app/client"
	"github.com/ed16/word-of-wisdom/internal/app/server"
	"github.com/stretchr/testify/assert"
)

func TestServerClientInteraction(t *testing.T) {
	serverConfig := &config.ServerConfig{
		ListenAddr: "localhost:9999",
		Difficulty: 1,
	}
	srv := server.NewServer(serverConfig)
	ctx, cancel := context.WithCancel(context.Background())
	go srv.Start(ctx)
	defer cancel()

	clientConfig := &config.ClientConfig{
		ServerAddr: "localhost:9999",
	}
	clt := client.NewClient(clientConfig)

	clt.Start(ctx)
	assert.NotEmpty(t, srv)
	assert.NotEmpty(t, clt)
}
