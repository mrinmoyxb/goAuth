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

func (s *UserService) ServiceSendEmailVerification(ctx context.Context, email string, token string) (bool, error) {
	if email == "" {
		return false, apperrors.ErrInvalidEmail
	}

	if token == "" {
		return false, apperrors.ErrInvalidEmailToken
	}

	// 14. send mail
	success, err := utils.SendMail(email, token)
	if success == false {
		fmt.Println("Error: ", err)
		return false, err
	}

	fmt.Println("Mail sent")
	return true, nil
}

func (s *UserService) ServiceEmailVerificationToken(ctx context.Context, email string, token string) (bool, error) {
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
	user := &models.UserEmailVerify{
		Email: email,
	}
	response := false
	err = s.repo.RepositoryFindByEmailEmailVerification(ctx, user, &response)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return false, apperrors.ErrInvalidCredentials
		}
		return false, err
	}

	// 2. verify email token
	tokenVerified, err := utils.UtilVerifyEntry(token, user.TokenHash)
	if err != nil || !tokenVerified {
		return false, apperrors.ErrInvalidCredentials
	}

	// 3. update users and email table
	err = s.repo.ReposistoryUpdateEmailVerification(ctx, user.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return false, apperrors.ErrInvalidCredentials
		}
		return false, err
	}

	return true, nil

}
