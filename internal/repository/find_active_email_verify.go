package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
)

func (r *PostgresRepository) FindActiveEmailVerification(ctx context.Context, verification *models.EmailVerification, response *bool) error {
	err := r.db.QueryRow(
		ctx,
		`SELECT token_hash, expires_at FROM email_verification WHERE email = $1 AND expires_at > current_timestamp`,
		verification.Email,
	).Scan(&verification.TokenHash, verification.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: verification record not found", apperrors.ErrUserNotFound)
		}
		return fmt.Errorf("%w: %v", apperrors.ErrQueryFailed, err)
	}
	fmt.Println("DB CHECK")
	return nil
}
