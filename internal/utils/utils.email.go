package utils

import (
	"goAuth/internal/apperrors"
	"net/mail"
	"strings"
)

func NormalizeEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if email == "" || len(email) > 254 {
		return "", apperrors.ErrInvalidEmail
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", apperrors.ErrInvalidEmail
	}
	return email, nil
}
