package repository

import (
	"context"
	"errors"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"

	"github.com/jackc/pgx/v5"
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

func (r *PostgresRepository) RepositoryCreateUser(ctx context.Context, user *models.UserInput, metadata *models.UserMetadata) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}
	defer tx.Rollback(ctx)
	fmt.Println("userID: ", user.ID)

	// 9. insert into users table
	fmt.Println("Check 1")
	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	fmt.Println("userID: ", user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrEmailAlreadyExists
		}
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	// 10. insert into metadata table
	fmt.Println("Check 2")
	metadata.UserID = user.ID
	_, err = tx.Exec(
		ctx,
		`INSERT INTO users_metadata (user_id, data) VALUES ($1, $2)`,
		metadata.UserID, metadata.Data,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}

	// 11. insert into email_verification table
	fmt.Println("userID: ", user.ID)
	fmt.Println("useremail: ", user.Email)
	fmt.Println("userToken: ", user.TokenHash)
	fmt.Println("Check 3")
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

	fmt.Println("Check 4")
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}
	return nil
}

func (r *PostgresRepository) RepositoryFindByEmail(ctx context.Context, email string) (*models.User, error) {
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

func (r *PostgresRepository) RepositoryFindByEmailEmailVerification(ctx context.Context, user *models.UserEmailVerify, response *bool) error {
	err := r.db.QueryRow(
		ctx,
		`SELECT token_hash, expires_at FROM email_verification WHERE email = $1 AND expires_at > current_timestamp`,
		user.Email,
	).Scan(&user.TokenHash, user.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: verification record not found", apperrors.ErrUserNotFound)
		}
		return fmt.Errorf("%w: %v", apperrors.ErrQueryFailed, err)
	}
	fmt.Println("DB CHECK")
	return nil
}

func (r *PostgresRepository) ReposistoryUpdateEmailVerification(ctx context.Context, email string, token string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrQueryFailed, err)
	}
	defer tx.Rollback(ctx)

	// update email verification table
	fmt.Println("DB1")
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
	fmt.Println("DB2")
	fmt.Println("userID: ", userID)

	// update users table
	tag, err := tx.Exec(
		ctx,
		`UPDATE users SET email_verified = true WHERE id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrQueryFailed, err)
	}
	fmt.Println("DB3")

	if tag.RowsAffected() != 1 {
		return errors.New("user associated with verification token not found")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrInsertFailed, err)
	}
	return nil
}
