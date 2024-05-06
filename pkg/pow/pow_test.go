package pow_test

import (
	"testing"

	"github.com/ed16/word-of-wisdom/pkg/pow"
	"github.com/stretchr/testify/assert"
)

func TestValidateChallengePositive(t *testing.T) {
	challenge := "testchallenge"
	difficulty := byte(4)
	nonce := pow.SolveChallenge(challenge, difficulty)
	valid := pow.ValidateChallenge(challenge, nonce, difficulty)
	assert.True(t, valid)
}

func TestValidateChallengeFalse(t *testing.T) {
	challenge := "testchallenge"
	difficulty := byte(4)
	nonce := "000001"
	valid := pow.ValidateChallenge(challenge, nonce, difficulty)
	assert.False(t, valid)
}
