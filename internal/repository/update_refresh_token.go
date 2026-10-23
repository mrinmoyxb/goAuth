package repository

import (
	"context"
	"goAuth/internal/apperrors"
)

func (r *PostgresRepository) UpdateRefreshToken(ctx context.Context, token string, userID int64) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE sessions SET refresh_token_hash = $1, expires_at = CURRENT_TIMESTAMP + interval '7 days' WHERE user_id = $2`,
		token, userID,
	)

	if err != nil {
		return apperrors.ErrQueryFailed
	}

	return nil
}