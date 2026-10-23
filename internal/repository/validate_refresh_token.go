package repository

import (
	"context"
	"errors"
	"goAuth/internal/apperrors"

	"github.com/jackc/pgx/v5"
)

type RefreshTokenInfo struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
}

func (r *PostgresRepository) ValidateRefreshToken(ctx context.Context, token string) (*RefreshTokenInfo, error) {
	refreshToken := &RefreshTokenInfo{}

	err := r.db.QueryRow(
		ctx,
		`SELECT id, user_id FROM sessions WHERE refresh_token_hash = $1 AND expires_at > CURRENT_TIMESTAMP`,
		token,
	).Scan(&refreshToken.ID, &refreshToken.UserID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrRefreshTokenInvalid
		}
		return nil, apperrors.ErrQueryFailed
	}

	return refreshToken, nil
}
