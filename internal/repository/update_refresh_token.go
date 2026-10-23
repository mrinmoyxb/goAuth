package repository

import (
	"context"
	"goAuth/internal/apperrors"
)

func (r *PostgresRepository) UpdateRefreshToken(ctx context.Context, oldToken string, newToken string, sessionID int64, userID int64) error {
	tag, err := r.db.Exec(
		ctx,
		`UPDATE sessions SET refresh_token_hash = $1, expires_at = CURRENT_TIMESTAMP + interval '7 days'
		WHERE id = $2 AND user_id = $3 AND refresh_token_hash = $4 AND expires_at > CURRENT_TIMESTAMP`,
		newToken, sessionID, userID, oldToken,
	)

	if err != nil {
		return apperrors.ErrQueryFailed
	}

	if tag.RowsAffected() != 1 {
		return apperrors.ErrRefreshTokenInvalid
	}

	return nil
}
