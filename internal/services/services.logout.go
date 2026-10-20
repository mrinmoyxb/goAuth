package services

import (
	"context"
	"goAuth/internal/apperrors"
	"goAuth/internal/utils"
)

func (s *UserService) ServiceLogout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return apperrors.ErrRefreshTokenInvalid
	}
	tokenHash := utils.HashVerificationToken(refreshToken)
	return s.repo.RepositoryRevokeCurrentSession(ctx, tokenHash)
}


func (s *UserService) ServiceLogoutAll(ctx context.Context, userID int64) error {
	return s.repo.RepositoryRevokeAllSession(ctx, userID)
}
