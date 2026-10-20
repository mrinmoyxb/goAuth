package repository

import (
	"context"
	"fmt"
	"goAuth/internal/apperrors"
	"errors"
	"github.com/jackc/pgx/v5"
)

type PasswordReset struct {
	UserID      int64
	Email       string
}

func (r *PostgresRepository) RepositoryPasswordReset(ctx context.Context, userID int64, email string, token string) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO password_reset (user_id, email, token_hash, expires_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP + interval '5 minutes')`,
		userID,
		email,
		token,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	return nil
}

func (r *PostgresRepository) FindByEmailasswordReset(ctx context.Context, email string, hashedToken string) (*PasswordReset, error) {
	res := &PasswordReset{}

	err := r.db.QueryRow(
		ctx,
		`SELECT user_id, email from password_reset WHERE email = $1 AND token_hash = $2 AND expires_at > CURRENT_TIMESTAMP AND used_at IS NULL`,
		email,
		hashedToken,
	).Scan(&res.UserID, &res.Email)
	if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, apperrors.ErrUserNotFound
        }

        return nil, fmt.Errorf("%w: %w", apperrors.ErrQueryFailed, err)
    }

	return res, nil
}

func (r *PostgresRepository) RepositoryResetPassword(ctx context.Context, userID int64, tokenHash string, newPassword string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}
	defer tx.Rollback(ctx)

	// 1. updated the password_reset table
	_, err = tx.Exec(
		ctx,
		`UPDATE password_reset SET used_at = CURRENT_TIMESTAMP WHERE token_hash = $1 AND user_id = $2 AND expires_at > CURRENT_TIMESTAMP AND used_at IS NULL`,
		tokenHash, userID,
	)
	if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return errors.New("invalid, expired, or already-used reset token")
        }
        return fmt.Errorf("%w: %w", apperrors.ErrQueryFailed, err)
    }

	// 2. update password in the users table
	tag , err := tx.Exec(
		ctx,
		`UPDATE users set password_hash = $1 AND updated_at = CURRENT_TIMESTAMP WHERE user_id = $2 AND email_verified = true`,
		newPassword, userID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrQueryFailed, err)
	}
	if tag.RowsAffected() != 1 {
        return errors.New("verified user not found or email not verified")
    }

	// 3. revoke refresh token
	_, err = tx.Exec(
		ctx,
		`UPDATE sessions SET revoked_at = CURRENT_TIMESTAMP WHERE user_id = $1 AND revoked IS NULL`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrQueryFailed, err)
	}
	
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	return nil
}
