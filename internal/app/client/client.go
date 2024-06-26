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
	config    *config.ClientConfig
	connector tcp.Connector
}

func NewClient(conf *config.ClientConfig, connector tcp.Connector) *Client {
	if connector == nil {
		connector = &tcp.DefaultConnector{}
	}
	return &Client{config: conf, connector: connector}
}

func (c *Client) Start(ctx context.Context) {
	for i := 0; i < c.config.RequestCount; i++ {
		if ctx.Err() != nil {
			break
		}
		quote, err := c.requestQuote(ctx)
		if err != nil {
			log.Printf("Failed to get quote: %s\n", err)
		}
		log.Printf("Received quote #%d: %s\n", i+1, quote)
	}
}

func (c *Client) requestQuote(ctx context.Context) (string, error) {
	conn, err := c.connector.Connect(ctx, c.config.ServerAddr)
	if err != nil {
		return "", err
	}
	defer c.connector.Close(conn)

	challengePrompt, err := c.connector.Receive(conn)
	if err != nil {
		return "", err
	}
	challenge, difficulty := extractChallengeAndDifficulty(challengePrompt)

	startTime := time.Now()
	nonce := pow.SolveChallenge(challenge, difficulty)
	log.Printf("Challenge solved in %v\n", time.Since(startTime))

	err = c.connector.Send(conn, nonce+"\n")
	if err != nil {
		return "", err
	}

	quote, err := c.connector.Receive(conn)
	return quote, err
}

func extractChallengeAndDifficulty(challengePrompt string) (string, byte) {
	var challenge string
	var difficulty int
	fmt.Sscanf(challengePrompt, "Solve PoW: SHA256( %s + <nonce> ) with %d leading zeros", &challenge, &difficulty)

	return challenge, byte(difficulty)
}
