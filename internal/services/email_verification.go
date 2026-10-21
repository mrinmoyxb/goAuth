package services

import (
	"context"
	"errors"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"goAuth/internal/utils"
	"net/mail"
	"strings"
)

func (s *UserService) SendVerificationEmail(ctx context.Context, email string, token string) (bool, error) {
	if email == "" {
		return false, apperrors.ErrInvalidEmail
	}

	if token == "" {
		return false, apperrors.ErrInvalidEmailToken
	}

	// 14. send mail
	success, err := utils.SendVerificationEmail(email, token)
	if success == false {
		fmt.Println("Error: ", err)
		return false, err
	}

	fmt.Println("Mail sent")
	return true, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, email string, token string) (bool, error) {
	if email == "" {
		return false, apperrors.ErrInvalidEmail
	}

	if token == "" {
		return false, apperrors.ErrInvalidEmailToken
	}

	email = strings.ToLower(strings.TrimSpace(email))
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, fmt.Errorf("%w: %v", apperrors.ErrInvalidEmail, err)
	}

	// 1. query email
	user := &models.EmailVerification{
		Email: email,
	}

	hashedToken := utils.HashToken(token)
	fmt.Println(hashedToken)

	// 3. update users and email table
	err = s.repo.MarkEmailVerifiedByToken(ctx, user.Email, hashedToken)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return false, apperrors.ErrInvalidCredentials
		}
		return false, err
	}

	return true, nil

}
