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

	hashedToken := utils.HashToken(refreshToken)

	userID, err := s.repo.ValidateRefreshToken(ctx, hashedToken)
	if err != nil {
		return nil, err
	}

	token, err := utils.NewTokenService()
	if err != nil {
		return nil, apperrors.ErrFailedToGenerateToken
	}

	user_id := strconv.FormatInt(userID, 10)
	user := utils.AuthenticatedUser{
		UserID: user_id,
		Role:   "user",
	}

	access_token, err := token.GenerateAccessToken(user)
	if err != nil {
		return nil, apperrors.ErrFailedToGenerateToken
	}

	newRefreshToken, err := utils.GenerateToken(32)
	newRefreshTokenHash := utils.HashToken(newRefreshToken)

	err = s.repo.

	return &AccessTokenResetResponse{
		AccessToken: access_token,
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
