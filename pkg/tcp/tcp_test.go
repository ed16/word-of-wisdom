package tcp_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/ed16/word-of-wisdom/pkg/tcp"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0") // Listen on a free port
	assert.NoError(t, err)
	defer listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		conn, _ := listener.Accept()
		conn.Close()
	}()

	_, err = tcp.Connect(ctx, listener.Addr().String())
	assert.NoError(t, err)
}

func TestSendReceive(t *testing.T) {
	server, client := net.Pipe()
	defer tcp.Close(server)
	defer tcp.Close(client)

	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Test Send
	go func() {
		err := tcp.Send(server, "hello")
		assert.NoError(t, err)
	}()

	// Test Receive
	msg, err := tcp.Receive(client)
	assert.NoError(t, err)
	assert.Equal(t, "hello", msg)
}
