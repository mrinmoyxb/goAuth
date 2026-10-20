package repository

import (
	"context"
	"fmt"
	"goAuth/internal/apperrors"
)

func (r *PostgresRepository) RepositoryRevokeCurrentSession(ctx context.Context, refreshToken string) error {
	_, err := r.db.Exec(
		ctx, 
		`UPDATE sessions SET revoked_at = CURRENT_TIMESTAMP WHERE refresh_token_hash = $1 AND revoked_at IS NULL`,
		refreshToken,
	)

	if err != nil {
        return fmt.Errorf("%w: %w", apperrors.ErrQueryFailed, err)
    }

    return nil
}

func (r *PostgresRepository) RepositoryRevokeAllSession(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE sessions SET revoked_at = CURRENT_TIMESTAMP WHERE user_id = $1 AND revoked_at IS NULL`,
		userID,
	)
	if err != nil {
        return fmt.Errorf("%w: %w", apperrors.ErrQueryFailed, err)
    }

    return nil
}