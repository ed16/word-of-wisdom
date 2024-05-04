package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// SolveChallenge solves the given Proof of Work challenge from the server.
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
