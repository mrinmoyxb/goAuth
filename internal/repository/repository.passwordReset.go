package repository

import (
	"context"
	"fmt"
	"goAuth/internal/apperrors"
)

type PasswordReset struct {
	UserID      int64
	Email       string
	TokenHashed string
}

func (r *PostgresRepository) RepositoryPasswordReset(ctx context.Context, userID int64, email string, token string) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO password_reset (user_id, email, token_hash, expires_at) VALUES ($1, $2, CURRENT_TIMESTAMP + interval '5 minutes')`,
		userID,
		email,
		token,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	return nil
}

func (r *PostgresRepository) FindByEmailasswordReset(ctx context.Context, email string) (*PasswordReset, error) {
	res := &PasswordReset{}

	err := r.db.QueryRow(
		ctx,
		`SELECT user_id, email, token_hash from password_reset WHERE email = $1`,
		email,
	).Scan(&res.UserID, &res.Email, &res.TokenHashed)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrQueryFailed, err)
	}

	return res, nil
}

func (r *PostgresRepository) RepositoryResetPassword(ctx context.Context, userID int64, newPassword string) {

}
