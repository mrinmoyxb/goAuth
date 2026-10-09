package utils

import (
	"errors"
	"os"
)

func GenerateRefreshToken() (string, error) {
	token := os.Getenv("REFRESH_TOKEN_SECRET_KEY")
	if len(token) < 64 {
		return "", errors.New("REFRESH_TOKEN_SECRET_KEY is less than 64")
	}

	return token, nil
}
