package repository

import (
	"context"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
)

func (r *PostgresRepository) RepositorySendEmailVerificationMail(ctx context.Context, user *models.UserEmailVerification) (*models.UserEmailVerificationResponse, error) {
	userResponse := &models.UserEmailVerificationResponse{}

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO email_verification (user_id, token_hash) VALUES ($1, $2) RETURNING id, email, expires_at`,
		user.UserID,
		user.TokenHash,
	).Scan(&userResponse.UserID, &userResponse.Email, &userResponse.ExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrQueryFailed, err)
	}

	return userResponse, nil
}
