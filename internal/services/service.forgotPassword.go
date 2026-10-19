package services

import (
	"context"
	"fmt"
	"goAuth/internal/apperrors"
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

func (s *UserService) ServiceValidateResetPasswordToken(ctx context.Context, email string, token string) error {
	// 1. verify email
	if email == "" {
		return apperrors.ErrInvalidEmail
	}
	email = strings.ToLower(strings.TrimSpace(email))
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidEmail, err)
	}

	// 2. hash token
	if token == "" {
		return apperrors.ErrInvalidToken
	}
	hashedToken := utils.HashVerificationToken(token)

	// 3.find the user by email and validating token
	existingUser, err := s.repo.FindByEmailasswordReset(ctx, email)
	if hashedToken != existingUser.TokenHashed {
		return apperrors.ErrInvalidPasswordResetToken
	}

	return nil

}
