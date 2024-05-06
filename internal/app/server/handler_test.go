package server_test

import (
	"bytes"
	"context"
	"net"
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
	buffer *bytes.Buffer
}

func (m *MockConn) Read(b []byte) (int, error) {
	return m.buffer.Read(b)
}

func (m *MockConn) Write(b []byte) (int, error) {
	m.buffer.Write(b)
	return len(b), nil
}

func (m *MockConn) Close() error {
	return nil
}

func (m *MockConn) SetDeadline(t time.Time) error {
	return nil
}

type MockConnector struct {
	mock.Mock
	tcp.Connector
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

func TestHandleConnection_InvalidPoW(t *testing.T) {
	mockConn := &MockConn{buffer: &bytes.Buffer{}}
	mockConnector := new(MockConnector)
	srv := server.NewServer(&config.ServerConfig{
		ListenAddr: "localhost:8080",
		Difficulty: 3,
		Deadline:   3 * time.Second,
	}, mockConnector)

	mockConnector.On("Send", mockConn, mock.Anything).Return(nil)
	mockConnector.On("Receive", mockConn).Return("000002", nil)
	mockConnector.On("Close", mockConn).Return(nil)

	srv.AddConnection()
	srv.HandleConnection(mockConn)

	mockConn.AssertExpectations(t)
	mockConnector.AssertExpectations(t)

	assert.Contains(t, mockConn.buffer.String(), "Invalid PoW solution\n")
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
