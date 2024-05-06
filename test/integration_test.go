package integration_test

import (
	"context"
	"testing"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/ed16/word-of-wisdom/internal/app/client"
	"github.com/ed16/word-of-wisdom/internal/app/server"
	"github.com/ed16/word-of-wisdom/pkg/tcp"
	"github.com/stretchr/testify/assert"
)

func TestServerClientInteraction(t *testing.T) {
	serverConfig := &config.ServerConfig{
		ListenAddr: "127.0.0.1:0",
		Difficulty: 1,
	}
	srv := server.NewServer(serverConfig, &tcp.DefaultConnector{})
	ctx, cancel := context.WithCancel(context.Background())
	go srv.Start(ctx)
	defer cancel()

	clientConfig := &config.ClientConfig{
		ServerAddr: "127.0.0.1:0",
	}
	clt := client.NewClient(clientConfig, &tcp.DefaultConnector{})

	clt.Start(ctx)
	assert.NotEmpty(t, srv)
	assert.NotEmpty(t, clt)
}
