package services

import (
	"context"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"goAuth/internal/repository"
	"goAuth/internal/utils"
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
	email, err := utils.NormalizeEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if !utils.ValidatePasswordPolicy(req.Password) {
		return nil, apperrors.ErrInvalidPassword
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	emailToken, err := utils.GenerateToken(10)
	if err != nil {
		return nil, fmt.Errorf("failed to generate email token: %w", err)
	}

	emailTokenHash := utils.HashToken(emailToken)
	if err != nil {
		return nil, fmt.Errorf("failed to hash email token: %w", err)
	}

	user := &models.UserInput{
		Email:        email,
		PasswordHash: passwordHash,
		TokenHash:    emailTokenHash,
	}

	data := req.Metadata
	if data == nil {
		data = map[string]any{}
	}

	metadata := &models.UserMetadata{
		Data: data,
	}

	if err := s.repo.CreateUser(ctx, user, metadata); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	result := &models.RegisterResponse{
		ID:            user.ID,
		Email:         user.Email,
		EmailVerified: false,
		EmailStatus:   "pending",
	}

	_, err = s.SendVerificationEmail(ctx, result.Email, emailToken)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	return result, nil
}
