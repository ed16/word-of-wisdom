package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/ed16/word-of-wisdom/internal/quotes"
	"github.com/ed16/word-of-wisdom/pkg/pow"
	"github.com/ed16/word-of-wisdom/pkg/tcp"
)

type Server struct {
	config    *config.ServerConfig
	connector tcp.Connector
	wg        sync.WaitGroup // WaitGroup to manage active connections
}

func NewServer(cfg *config.ServerConfig, connector tcp.Connector) *Server {
	if connector == nil {
		connector = &tcp.DefaultConnector{}
	}
	return &Server{config: cfg, connector: connector}
}

// Start begins listening and accepting connections on the server's address.
func (s *Server) Start(ctx context.Context) {
	lc := net.ListenConfig{
		KeepAlive: s.config.KeepAlive,
	}

	listener, err := lc.Listen(ctx, "tcp", s.config.ListenAddr)

	if err != nil {
		log.Panicf("Failed to listen on %s: %v", s.config.ListenAddr, err)
	}
	defer listener.Close()

	log.Printf("Server started, listening on %s\n", s.config.ListenAddr)

	go func() {
		<-ctx.Done()
		log.Println("Gracefully stopping...")
		s.wg.Wait() // Wait for all connections to finish
		log.Println("All connections closed.")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Check if the error is due to context cancellation.
			if ctx.Err() != nil {
				log.Println("Server shutting down...")
				return
			}
			log.Println("Error accepting connection:", err)
			continue
		}
		s.AddConnection() // Increment WaitGroup counter
		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer func() {
		s.connector.Close(conn)
		s.DoneConnection() // Decrement WaitGroup counter on connection close
	}()
	_ = conn.SetDeadline(time.Now().Add(s.config.Deadline))

	challenge := issueChallenge()
	// Send challenge
	err := s.connector.Send(conn, fmt.Sprintf("Solve PoW: SHA256( %s + <nonce> ) with %d leading zeros\n", challenge, s.config.Difficulty))
	if err != nil {
		log.Println("Error sending challenge:", err)
	}

	// Receive solution
	solution, err := s.connector.Receive(conn)
	if err != nil {
		log.Println("Error receiving solution:", err)
	}

	if pow.ValidateChallenge(challenge, solution, s.config.Difficulty) {
		quote := quotes.GetRandomQuote()
		_, err = conn.Write([]byte(quote + "\n"))
		if err != nil {
			log.Println("Error sending quote:", err)
		}
	} else {
		_, err = conn.Write([]byte("Invalid PoW solution\n"))
		if err != nil {
			log.Println("Error sending responce:", err)
		}
	}
}

func (s *Server) AddConnection() {
	s.wg.Add(1)
}

func (s *Server) DoneConnection() {
	s.wg.Done()
}

func issueChallenge() string {
	challenge := strconv.Itoa(rand.Intn(1000000))
	return challenge
}
