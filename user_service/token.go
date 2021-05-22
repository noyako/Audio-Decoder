package main

import (
	"crypto/rand"
	"encoding/hex"
)

const tokenLength = 30

func generateToken() string {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
