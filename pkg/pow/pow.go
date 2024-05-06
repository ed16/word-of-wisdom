package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// IssueChallenge generates the Proof of Work challenge
func IssueChallenge() string {
	// Current time as a seed to ensure time-based uniqueness
	timeSeed := []byte(time.Now().UTC().Format(time.RFC3339Nano))

	// Generate a cryptographically secure random number
	randomBytes := make([]byte, 16) // 128 bits
	if _, err := rand.Read(randomBytes); err != nil {
		log.Panic("Failed to generate random bytes: ", err)
	}

	// Combine time seed and random bytes
	challengeBytes := append(timeSeed, randomBytes...)

	// Hash the combined bytes to get a fixed size output
	hash := sha256.Sum256(challengeBytes)

	// Convert the hash to a hexadecimal string
	challenge := hex.EncodeToString(hash[:])

	return challenge
}

// SolveChallenge solves the given Proof of Work challenge from the server
func SolveChallenge(challenge string, difficulty byte) string {
	// Brute force search for the nonce
	var nonce int
	for {
		nonceStr := strconv.Itoa(nonce)
		if ValidateChallenge(challenge, nonceStr, difficulty) {
			return nonceStr
		}
		nonce++
	}
}

func ValidateChallenge(challenge, nonce string, difficulty byte) bool {
	// Concatenate challenge and nonce
	data := fmt.Sprintf("%s%s", challenge, nonce)
	// Compute SHA-256 hash
	hash := sha256.Sum256([]byte(data))
	hexHash := hex.EncodeToString(hash[:])
	// Check if hash has the required number of leading zeros
	return strings.HasPrefix(hexHash, strings.Repeat("0", int(difficulty)))
}
