package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomToken(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, n)
	for i := range n {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		token[i] = letters[num.Int64()]
	}
	return string(token), nil
}
