package repository

import (
	"context"
	"fmt"
	"goAuth/internal/apperrors"
)

func (r *PostgresRepository) CreateSession(ctx context.Context, userID int64, refreshToken string) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO sessions (user_id, refresh_token_hash, expires_at) VALUES ($1, $2, CURRENT_TIMESTAMP + interval '15 days')`,
		userID, refreshToken,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	return nil
}
