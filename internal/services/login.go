package services

import (
	"context"
	"errors"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"goAuth/internal/utils"
	"net/mail"
	"strconv"
	"strings"
)

func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	if req.Email == "" {
		return nil, apperrors.ErrInvalidEmail
	}

	if req.Password == "" {
		return nil, apperrors.ErrInvalidPassword
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidEmail, err)
	}

	passwordVerified := utils.ValidatePasswordPolicy(req.Password)
	if !passwordVerified {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidPassword, err)
	}

	existingUser, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if existingUser.EmailVerified == false {
		return nil, apperrors.ErrEmailNotVerified
	}

	passVerified, err := utils.VerifyPassword(req.Password, existingUser.PasswordHash)
	if err != nil || !passVerified {
		return nil, apperrors.ErrInvalidCredentials
	}

	userToken := utils.AuthenticatedUser{UserID: strconv.FormatInt(existingUser.ID, 10), Role: "user"}
	token, err := utils.NewTokenService()
	if err != nil {
		fmt.Println("error")
	}
	accessToken, err := token.GenerateAccessToken(userToken)
	if err != nil {
		fmt.Println("error")
	}

	refreshToken, err := utils.GenerateToken(32)
	refreshTokenHash := utils.HashToken(refreshToken)

	err = s.repo.CreateSession(ctx, existingUser.ID, refreshTokenHash)
	if err != nil {
		fmt.Println("Failed to insert refresh token in db")
		return nil, err
	}

	return &models.LoginResponse{
		Message:      "login successfull",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
