package repository

import (
	"context"
	"goAuth/internal/apperrors"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) ValidateRefreshToken(ctx context.Context, token string) (int64, error) {
	var userID int64
	err := r.db.QueryRow(
		ctx,
		`SELECT user_id FROM session WHERE refresh_token_hash = $1 AND expires_at > CURRENT_TIMESTAMP`,
		token,
	).Scan(&userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, apperrors.ErrRefreshTokenInvalid
		}
		return -1, apperrors.ErrQueryFailed
	}

	return userID, nil
}