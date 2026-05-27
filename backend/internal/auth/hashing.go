package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

type sha256TokenHasher struct{}

func (h *sha256TokenHasher) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])
	return tokenHash
}
