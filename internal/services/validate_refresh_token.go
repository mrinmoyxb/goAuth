package services

import (
	"context"
	"goAuth/internal/apperrors"
	"goAuth/internal/utils"
	"strconv"
)

type AccessTokenResetResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *UserService) RefreshAccessToken(ctx context.Context, refreshToken string) (*AccessTokenResetResponse, error) {
	if refreshToken == "" {
		return nil, apperrors.ErrRefreshTokenInvalid
	}
	oldHashedToken := utils.HashToken(refreshToken)

	refreshTokenInfo, err := s.repo.ValidateRefreshToken(ctx, oldHashedToken)
	if err != nil {
		return nil, err
	}

	token, err := utils.NewTokenService()
	if err != nil {
		return nil, apperrors.ErrFailedToGenerateToken
	}
	user := utils.AuthenticatedUser{
		UserID: strconv.FormatInt(refreshTokenInfo.UserID, 10),
		Role:   "user",
	}
	access_token, err := token.GenerateAccessToken(user)
	if err != nil {
		return nil, apperrors.ErrFailedToGenerateToken
	}

	newRefreshToken, err := utils.GenerateToken(32)
	if err != nil {
		return nil, apperrors.ErrFailedToGenerateToken
	}
	newRefreshTokenHash := utils.HashToken(newRefreshToken)

	err = s.repo.UpdateRefreshToken(ctx, oldHashedToken, newRefreshTokenHash, refreshTokenInfo.ID, refreshTokenInfo.UserID)
	if err != nil {
		return nil, err
	}

	return &AccessTokenResetResponse{
		AccessToken:  access_token,
		RefreshToken: newRefreshToken,
	}, nil
}

/*
Your current flow issues a new access token but leaves the existing refresh token valid. For a production-oriented authentication system, consider rotating refresh tokens:
1. Validate the existing refresh token.
2. Generate a new access token and refresh token.
3. Replace the old refresh-token hash in the database atomically.
4. Send the new refresh token in an HttpOnly cookie.

refresh-token rotation, session revocation checks, and appropriate user-status validation.
*/
