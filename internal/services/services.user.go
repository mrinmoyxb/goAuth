package services

import (
	"context"
	"errors"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"goAuth/internal/repository"
	"goAuth/internal/utils"
	"strings"
)

type UserService struct {
	repo repository.PostgresRepository
}

func NewUserService(repo repository.PostgresRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) ServiceCreateUser(ctx context.Context, req models.UserRequest) (*models.UserResponse, error) {
	if req.Email == "" {
		return nil, apperrors.ErrInvalidEmail
	}

	if req.Password == "" {
		return nil, apperrors.ErrInvalidPassword
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, err := s.repo.RepositoryFindByEmail(ctx, email)
	if err != nil && !errors.Is(err, apperrors.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, apperrors.ErrEmailAlreadyExists
	}

	passwordHash, err := utils.UtilHashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	metadata := &models.UserMetadata{
		Data: req.Metadata,
	}

	if err := s.repo.RepositoryCreateUser(ctx, user, metadata); err != nil {
		if errors.Is(err, apperrors.ErrEmailAlreadyExists) {
			return nil, apperrors.ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Metadata:  metadata.Data,
	}, nil
}

func (s *UserService) ServiceLogin(ctx context.Context, req models.UserRequestLogin) (*models.UserResponse, error) {
	if req.Email == "" {
		return nil, apperrors.ErrInvalidEmail
	}

	if req.Password == "" {
		return nil, apperrors.ErrInvalidPassword
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, err := s.repo.RepositoryFindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	passVerified, err := utils.UtilVerifyPassword(req.Password, existingUser.PasswordHash)
	if err != nil || !passVerified {
		return nil, apperrors.ErrInvalidCredentials
	}

	return &models.UserResponse{
		ID:        existingUser.ID,
		Email:     existingUser.Email,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt,
	}, nil
}
