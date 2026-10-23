package services

import (
	"context"
	"goAuth/internal/apperrors"
	"goAuth/internal/utils"
)

func (s *UserService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return apperrors.ErrRefreshTokenInvalid
	}
	tokenHash := utils.HashToken(refreshToken)
	return s.repo.RevokeSessionByToken(ctx, tokenHash)
}


