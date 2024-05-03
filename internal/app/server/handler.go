package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/ed16/word-of-wisdom/internal/quotes"
	"github.com/ed16/word-of-wisdom/pkg/pow"
	"github.com/ed16/word-of-wisdom/pkg/tcp"
)

type Server struct {
	Config *config.ServerConfig
}

func NewServer(cfg *config.ServerConfig) *Server {
	return &Server{Config: cfg}
}

// Start begins listening and accepting connections on the server's address.
func (s *Server) Start(ctx context.Context) {
	lc := net.ListenConfig{
		KeepAlive: s.Config.KeepAlive,
	}

	listener, err := lc.Listen(ctx, "tcp", s.Config.ListenAddr)

	if err != nil {
		log.Panicf("Failed to listen on %s: %v", s.Config.ListenAddr, err)
	}
	defer listener.Close()

	log.Printf("Server started, listening on %s\n", s.Config.ListenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	challenge := issueChallenge()
	// Send challenge
	err := tcp.Send(conn, fmt.Sprintf("Solve PoW: SHA256( %s + <nonce> ) with %d leading zeros\n", challenge, s.Config.Difficulty))
	if err != nil {
		log.Println("Error sending challenge:", err)
	}

	// Read response
	response, err := tcp.Receive(conn)
	if err != nil {
		log.Println("Error receiving response:", err)
	}

	if pow.ValidateChallenge(challenge, response, s.Config.Difficulty) {
		quote := quotes.GetRandomQuote()
		conn.Write([]byte(quote + "\n"))
	} else {
		conn.Write([]byte("Invalid PoW solution\n"))
	}
}

func issueChallenge() string {
	challenge := strconv.Itoa(rand.Intn(1000000))
	return challenge
}
