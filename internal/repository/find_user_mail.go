package repository

import (
	"context"
	"errors"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	err := r.db.QueryRow(
		ctx,
		`SELECT id, email, password_hash, created_at, updated_at,  FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.EmailVerified,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", apperrors.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", apperrors.ErrQueryFailed, err)
	}

	return user, nil
}
