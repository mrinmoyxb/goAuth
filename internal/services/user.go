package services

import (
	"context"
	"errors"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"goAuth/internal/repository"
	"goAuth/internal/utils"
	"net/mail"
	"strconv"
	"strings"
)

type UserService struct {
	repo *repository.PostgresRepository
}

func NewUserService(repo *repository.PostgresRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, req models.RegisterRequest) (*models.RegisterResponse, error) {
	// 1. normalized mail
	email, err := utils.NormalizeEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// 2. password validation
	if !utils.ValidatePasswordPolicy(req.Password) {
		return nil, apperrors.ErrInvalidPassword
	}

	// 3. hashed password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. email verification token
	emailToken, err := utils.GenerateToken(10)
	if err != nil {
		return nil, fmt.Errorf("failed to generate email token: %w", err)
	}

	// 5. hashed email verification token
	emailTokenHash := utils.HashToken(emailToken)
	if err != nil {
		return nil, fmt.Errorf("failed to hash email token: %w", err)
	}

	// 6. need to insert in the users table
	user := &models.UserInput{
		Email:        email,
		PasswordHash: passwordHash,
		TokenHash:    emailTokenHash,
	}

	// 7. metadata
	data := req.Metadata
	if data == nil {
		data = map[string]any{}
	}

	metadata := &models.UserMetadata{
		Data: data,
	}

	// 8. call to db to create user
	if err := s.repo.CreateUser(ctx, user, metadata); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// 12. response from DB
	result := &models.RegisterResponse{
		ID:               user.ID,
		Email:            user.Email,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		Metadata:         metadata.Data,
		VerificationMail: false,
	}

	// 13. call service for email verification
	success, err := s.SendVerificationEmail(ctx, result.Email, emailToken)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}
	if success == true {
		result.VerificationMail = true
	}

	// 15. done
	return result, nil
}

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
