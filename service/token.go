package service

import (
	"crypto/rand"
	"encoding/hex"
)

const tokenLength = 30

// GenerateToken generates random string
func GenerateToken() string {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
