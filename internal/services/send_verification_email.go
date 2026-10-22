package services

import (
	"context"
	"log"
	"goAuth/internal/apperrors"
	"goAuth/internal/utils"
)

func (s *UserService) SendVerificationEmail(ctx context.Context, email string, token string) (bool, error) {
	if email == "" {
		return false, apperrors.ErrInvalidEmail
	}

	if token == "" {
		return false, apperrors.ErrInvalidEmailToken
	}

	go func(email, token string) (bool, error) {
	_, err := utils.SendVerificationEmail(email, token)
		if err != nil {
			log.Printf("send verification email: %v", err)
		}
		return false, err
	}(email, token)

	return true, nil
}