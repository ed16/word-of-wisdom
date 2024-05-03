package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/ed16/word-of-wisdom/pkg/pow"
	"github.com/ed16/word-of-wisdom/pkg/tcp"
)

type Client struct {
	conf *config.ClientConfig
}

func NewClient(conf *config.ClientConfig) *Client {
	return &Client{conf: conf}
}

func (c *Client) RequestQuote(ctx context.Context) (string, error) {
	conn, err := tcp.Connect(ctx, c.conf.ServerAddr)
	if err != nil {
		return "", err
	}
	defer tcp.Close(conn)

	challengePrompt, err := tcp.Receive(conn)
	if err != nil {
		return "", err
	}
	challenge, difficulty := extractChallengeAndDifficulty(challengePrompt)

	startTime := time.Now()
	nonce := pow.SolveChallenge(challenge, difficulty)
	log.Printf("Challenge solved in %v\n", time.Since(startTime))

	err = tcp.Send(conn, nonce+"\n")
	if err != nil {
		return "", err
	}

	quote, err := tcp.Receive(conn)
	return quote, err
}

func extractChallengeAndDifficulty(challengePrompt string) (string, byte) {
	var challenge string
	var difficulty int
	fmt.Sscanf(challengePrompt, "Solve PoW: SHA256( %s + <nonce> ) with %d leading zeros\n", &challenge, &difficulty)

	return challenge, byte(difficulty)
}