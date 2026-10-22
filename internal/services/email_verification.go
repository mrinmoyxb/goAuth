package services

import (
	"context"
	"errors"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"goAuth/internal/utils"
	"net/mail"
	"strings"
)

func (s *UserService) VerifyEmail(ctx context.Context, email string, token string) error {
	if email == "" {
		return apperrors.ErrInvalidEmail
	}

	if token == "" {
		return apperrors.ErrInvalidEmailToken
	}

	email = strings.ToLower(strings.TrimSpace(email))
	_, err := mail.ParseAddress(email)
	if err != nil {
		return apperrors.ErrInvalidEmail
	}

	user := &models.EmailVerification{
		Email: email,
	}

	hashedToken := utils.HashToken(token)

	err = s.repo.MarkEmailVerifiedByToken(ctx, user.Email, hashedToken)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrInvalidCredentials
		}
		return err
	}

	return nil

}
