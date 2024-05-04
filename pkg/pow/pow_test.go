package pow_test

import (
	"testing"

	"github.com/ed16/word-of-wisdom/pkg/pow"
	"github.com/stretchr/testify/assert"
)

func TestValidateChallenge(t *testing.T) {
	challenge := "testchallenge"
	difficulty := byte(4)
	nonce := pow.SolveChallenge(challenge, difficulty)
	valid := pow.ValidateChallenge(challenge, nonce, difficulty)
	assert.True(t, valid)
}
