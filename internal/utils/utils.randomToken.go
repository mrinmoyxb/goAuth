package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func GenerateRandomToken(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("token length must be positive")
	}
	raw := make([]byte, n)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}

func HashVerificationToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
