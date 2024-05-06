package pow_test

import (
	"testing"

	"github.com/ed16/word-of-wisdom/pkg/pow"
	"github.com/stretchr/testify/assert"
)

func TestIssueChallenge(t *testing.T) {
	// Create a set to hold challenges to test for uniqueness
	challenges := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		challenge := pow.IssueChallenge()
		if len(challenge) != 64 { // SHA-256 hash is 64 hex characters
			t.Errorf("Expected challenge length of 64, got %d", len(challenge))
		}
		if _, exists := challenges[challenge]; exists {
			t.Errorf("Duplicate challenge found: %s", challenge)
		}
		challenges[challenge] = true
	}
}

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
