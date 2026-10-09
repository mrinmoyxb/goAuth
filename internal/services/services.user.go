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

func (s *UserService) ServiceCreateUser(ctx context.Context, req models.UserRequest) (*models.UserResponse, error) {
	email, err := utils.NormalizeEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if !utils.ValidatePassword(req.Password) {
		return nil, apperrors.ErrInvalidPassword
	}

	passwordHash, err := utils.UtilHashEntry(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	data := req.Metadata
	if data == nil {
		data = map[string]any{}
	}

	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	metadata := &models.UserMetadata{
		Data: data,
	}

	if err := s.repo.RepositoryCreateUser(ctx, user, metadata); err != nil {
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
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidEmail, err)
	}

	passwordVerified := utils.ValidatePassword(req.Password)
	if !passwordVerified {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidPassword, err)
	}

	existingUser, err := s.repo.RepositoryFindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	passVerified, err := utils.UtilVerifyEntry(req.Password, existingUser.PasswordHash)
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
