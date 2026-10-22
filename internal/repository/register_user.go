package repository

import (
	"context"
	"errors"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.UserInput, metadata *models.UserMetadata) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrEmailAlreadyExists
		}
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	metadata.UserID = user.ID
	_, err = tx.Exec(
		ctx,
		`INSERT INTO users_metadata (user_id, data) VALUES ($1, $2)`,
		metadata.UserID, metadata.Data,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	_, err = tx.Exec(
		ctx,
		`INSERT into email_verification (user_id, email, token_hash, expires_at) VALUES ($1, $2, $3, current_timestamp + interval '5 minutes')`,
		user.ID,
		user.Email,
		user.TokenHash,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}
	return nil
}

