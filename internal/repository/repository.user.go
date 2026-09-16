package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"goAuth/internal/models"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
		`,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}
