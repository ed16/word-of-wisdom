package tcp

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
)

// Connector defines an interface for connecting and communicating over TCP.
type Connector interface {
	Connect(ctx context.Context, address string) (net.Conn, error)
	Send(conn net.Conn, message string) error
	Receive(conn net.Conn) (string, error)
	Close(conn net.Conn)
}

// DefaultConnector provides a default implementation of the Connector interface.
type DefaultConnector struct{}

func (d *DefaultConnector) Connect(ctx context.Context, address string) (net.Conn, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to %s: %w", address, err)
	}
	return conn, nil
}

func (d *DefaultConnector) Send(conn net.Conn, message string) error {
	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString(message + "\n"); err != nil {
		return fmt.Errorf("Failed to write message to connection: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("Failed to flush connection: %w", err)
	}
	return nil
}

func (d *DefaultConnector) Receive(conn net.Conn) (string, error) {
	const maxBufferSize = 16 * 1024 // 16KB
	reader := bufio.NewReaderSize(conn, maxBufferSize)
	message, err := reader.ReadString('\n')
	if len(message) > 0 {
		message = message[:len(message)-1]
	}
	if err != nil {
		return "", fmt.Errorf("Failed to read from connection: %w", err)
	}
	return message, nil
}

func (d *DefaultConnector) Close(conn net.Conn) {
	if err := conn.Close(); err != nil {
		log.Println("Failed to close connection:", err)
	}
}
