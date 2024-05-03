package tcp

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
)

// Connect establishes a TCP connection to the specified address.
func Connect(ctx context.Context, address string) (net.Conn, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	return conn, nil
}

// Send writes the given message to the TCP connection and flushes it.
func Send(conn net.Conn, message string) error {
	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString(message + "\n"); err != nil {
		return fmt.Errorf("failed to write message to connection: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush connection: %w", err)
	}
	return nil
}

// Receive reads a line from the TCP connection.
func Receive(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read from connection: %w", err)
	}
	return message, nil
}

// Close safely closes the TCP connection.
func Close(conn net.Conn) {
	if err := conn.Close(); err != nil {
		log.Println("Failed to close connection:", err)
	}
}
