package services

import (
	"context"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/repository"
	"goAuth/internal/utils"
	"net/mail"
	"strings"
)

// 1. user clicks on forgot password, provide the mail and auth system sends a token
func (s *UserService) ServiceForgotPasword(ctx context.Context, email string) error {
	// 1. check email and normalize it
	if email == "" {
		return apperrors.ErrInvalidEmail
	}
	email = strings.ToLower(strings.TrimSpace(email))
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidEmail, err)
	}

	// 2. check if that email exists
	existingUser, err := s.repo.RepositoryFindByEmail(ctx, email)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	// 3. generated token and hashed it
	token, err := utils.GenerateRandomToken(32)
	if err != nil {
		return apperrors.ErrFailedToGenerateToken
	}
	hashedToken := utils.HashVerificationToken(token)

	// 4. create a field/entry in DB -> email, token and expiry time
	err = s.repo.RepositoryPasswordReset(ctx, existingUser.ID, email, hashedToken)
	if err != nil {
		return apperrors.ErrInsertFailed
	}

	// 5. send token to the user in mail
	_, err = utils.SendMailForgotPassword(email, token)
	if err != nil {
		return apperrors.ErrFailedToSendMail
	}

	return nil
}

func (s *UserService) ServiceValidateResetPasswordToken(ctx context.Context, email string, token string) (*repository.PasswordReset, error) {
	// 1. verify email
	if email == "" {
		return nil, apperrors.ErrInvalidEmail
	}
	email = strings.ToLower(strings.TrimSpace(email))
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidEmail, err)
	}

	// 2. hash token
	if token == "" {
		return nil, apperrors.ErrInvalidToken
	}
	hashedToken := utils.HashVerificationToken(token)

	// 3.find the user by email and validate token
	existingUser, err := s.repo.FindByEmailasswordReset(ctx, email, hashedToken)
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	return existingUser, nil
}

func (s *UserService) ServiceValidateNewPassword(ctx context.Context, userID int64, token string, password string) error {

	// 1. hash token
	if token == "" {
		return apperrors.ErrInvalidToken
	}
	hashedToken := utils.HashVerificationToken(token)

	// 2. validate new password and hash it
	if !utils.ValidatePassword(password) {
		return apperrors.ErrInvalidPassword
	}
	passwordHash, err := utils.UtilHashEntry(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.repo.RepositoryResetPassword(ctx, userID, hashedToken, passwordHash)
	if err != nil {
		return apperrors.ErrInvalidCredentials
	}

	return nil
}
