package repository

import (
	"context"
	"errors"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"

	"github.com/jackc/pgx/v5"
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

func (r *PostgresRepository) RepositoryCreateUser(ctx context.Context, user *models.User, metadata *models.UserMetadata) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperrors.ErrTxBeginFailed
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return apperrors.ErrInsertFailed
	}

	metadata.UserID = user.ID
	_, err = tx.Exec(
		ctx,
		`INSERT INTO users_metadata (user_id, data) VALUES ($1 '{}'::jsonb)`,
		metadata.UserID, metadata.Data,
	)
	if err != nil {
		return apperrors.ErrInsertFailed
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) RepositoryFindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	err := r.db.QueryRow(
		ctx,
		`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w %v", apperrors.ErrQueryFailed, err)
	}

	return user, nil
}


