package quotes_test

import (
	"testing"

	"github.com/ed16/word-of-wisdom/internal/quotes"
	"github.com/stretchr/testify/assert"
)

func TestGetRandomQuote(t *testing.T) {
	quote := quotes.GetRandomQuote()
	assert.NotEmpty(t, quote)
}
