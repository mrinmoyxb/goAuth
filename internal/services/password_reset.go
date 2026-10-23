package services

import (
	"context"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/utils"
	"net/mail"
	"strings"
)

type ResponsePasswordReset struct {
	UserID int64 `json:"userID"`
}

func (s *UserService) RequestPasswordReset(ctx context.Context, email string) (*ResponsePasswordReset, error) {
	user := &ResponsePasswordReset{}

	if email == "" {
		return nil, apperrors.ErrInvalidEmail
	}
	email = strings.ToLower(strings.TrimSpace(email))
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidEmail, err)
	}

	fmt.Println("Check2")
	existingUser, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}
	user.UserID = existingUser.ID

	fmt.Println("Check3")
	token, err := utils.GenerateToken(32)
	if err != nil {
		return nil, apperrors.ErrFailedToGenerateToken
	}
	hashedToken := utils.HashToken(token)

	fmt.Println("Check4")
	err = s.repo.CreatePasswordReset(ctx, existingUser.ID, email, hashedToken)
	if err != nil {
		return nil, apperrors.ErrInsertFailed
	}

	fmt.Println("Check5")
	go func() error {
		_, err = utils.SendPasswordResetEmail(email, token)
		if err != nil {
			return apperrors.ErrFailedToSendMail
		}
		return nil
	}()
		
	return user, nil
}

func (s *UserService) ValidatePasswordResetToken(ctx context.Context, userID int64, token string) error {
	if token == "" {
		return apperrors.ErrInvalidToken
	}
	hashedToken := utils.HashToken(token)

	err := s.repo.FindPasswordReset(ctx, userID, hashedToken)
	if err != nil {
		return apperrors.ErrInvalidCredentials
	}

	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, userID int64, token string, password string) error {

	// 1. hash token
	if token == "" {
		return apperrors.ErrInvalidToken
	}
	hashedToken := utils.HashToken(token)

	// 2. validate new password and hash it
	if !utils.ValidatePasswordPolicy(password) {
		return apperrors.ErrInvalidPassword
	}
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.repo.ResetPassword(ctx, userID, hashedToken, passwordHash)
	if err != nil {
		return apperrors.ErrInvalidCredentials
	}

	return nil
}
