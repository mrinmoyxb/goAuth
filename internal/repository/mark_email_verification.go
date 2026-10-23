package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"goAuth/internal/apperrors"
)

func (r *PostgresRepository) MarkEmailVerifiedByToken(ctx context.Context, email string, token string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrQueryFailed, err)
	}
	defer tx.Rollback(ctx)

	var userID int64
	err = tx.QueryRow(
		ctx,
		`UPDATE email_verification SET used_at = CURRENT_TIMESTAMP WHERE token_hash = $1 AND expires_at > CURRENT_TIMESTAMP AND used_at IS NULL RETURNING user_id`,
		token,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("invalid, expired, or already-used verification token")
		}
		return fmt.Errorf("%w: %w", apperrors.ErrQueryFailed, err)
	}

	tag, err := tx.Exec(
		ctx,
		`UPDATE users SET email_verified = true WHERE id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrQueryFailed, err)
	}

	if tag.RowsAffected() != 1 {
		return errors.New("user associated with verification token not found")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	return nil
}
