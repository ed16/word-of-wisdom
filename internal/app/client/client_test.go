package client_test

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/ed16/word-of-wisdom/internal/app/client"
	"github.com/ed16/word-of-wisdom/internal/app/server"
	"github.com/ed16/word-of-wisdom/pkg/tcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConn struct {
	mock.Mock
	net.Conn
}

func (m *MockConn) Read(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Write(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Close() error {
	return m.Called().Error(0)
}

type MockConnector struct {
	mock.Mock
}

func (m *MockConnector) Connect(ctx context.Context, address string) (net.Conn, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(net.Conn), args.Error(1)
}

func (m *MockConnector) Send(conn net.Conn, message string) error {
	return m.Called(conn, message).Error(0)
}

func (m *MockConnector) Receive(conn net.Conn) (string, error) {
	args := m.Called(conn)
	return args.String(0), args.Error(1)
}

func (m *MockConnector) Close(conn net.Conn) {
	m.Called(conn)
}

func TestClientStart(t *testing.T) {
	cfg := &config.ClientConfig{
		ServerAddr:   "localhost:8080",
		RequestCount: 1,
	}
	mockConnector := new(MockConnector)
	mockConn := new(MockConn)
	clt := client.NewClient(cfg, mockConnector)

	mockConnector.On("Connect", mock.Anything, cfg.ServerAddr).Return(mockConn, nil)
	mockConnector.On("Send", mockConn, mock.Anything).Return(nil)
	mockConnector.On("Receive", mockConn).Return("Solve PoW: SHA256( challenge + <nonce> ) with 2 leading zeros\n", nil).Once()
	mockConnector.On("Receive", mockConn).Return("Here is your quote.", nil).Once()
	mockConnector.On("Close", mockConn).Return(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	clt.Start(ctx)

	mockConnector.AssertExpectations(t)
}

func TestConcurrentClientConnections(t *testing.T) {
	cfg := &config.ClientConfig{
		ServerAddr:   "localhost:8080",
		RequestCount: 1,
	}
	wg := sync.WaitGroup{}
	const numClients = 5
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			mockConnector := new(MockConnector)
			mockConn := new(MockConn)
			clt := client.NewClient(cfg, mockConnector)

			// Set up the mocks with expected calls for each client
			mockConnector.On("Connect", mock.Anything, cfg.ServerAddr).Return(mockConn, nil)
			mockConnector.On("Send", mockConn, mock.Anything).Return(nil)
			mockConnector.On("Receive", mockConn).Return("Here is your quote.", nil)
			mockConnector.On("Close", mockConn).Return(nil)

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			clt.Start(ctx)

			mockConnector.AssertExpectations(t)
		}(i)
	}
	wg.Wait()
}

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
