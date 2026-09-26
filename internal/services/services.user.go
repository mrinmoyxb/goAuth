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

func (s *UserService) SerivceCreateUser(ctx context.Context, req models.UserRequest) (*models.UserResponse, error) {
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

	passwordHash, err := utils.utilHashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email: email,
		PasswordHash: passwordHash,
	}
	metadata := &models.UserMetadata{
		Data: req.Metadata, 
	}

	return nil, nil
}
