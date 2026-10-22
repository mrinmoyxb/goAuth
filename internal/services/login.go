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

	// 1. Find user by email
	existingUser, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. Ensure email is verified
	if existingUser.EmailVerified == false {
		return nil, apperrors.ErrEmailNotVerified
	}

	// 3. Verify password
	passVerified, err := utils.VerifyPassword(req.Password, existingUser.PasswordHash)
	if err != nil || !passVerified {
		return nil, apperrors.ErrInvalidCredentials
	}

	// 4. Access token
	userToken := utils.AuthenticatedUser{UserID: strconv.FormatInt(existingUser.ID, 10), Role: "user"}
	token, err := utils.NewTokenService()
	if err != nil {
		fmt.Println("error")
	}
	accessToken, err := token.GenerateAccessToken(userToken)
	if err != nil {
		fmt.Println("error")
	}

	// 5. Refresh token
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