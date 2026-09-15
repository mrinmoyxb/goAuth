package database

import (
	"context"
	"goAuth/internal/apperrors"
	"os"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres() (*pgxpool.Pool, error){
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, apperrors.ErrDBURLNotFound;
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, apperrors.ErrParseDBURL;
	}

	config.MaxConns = 10
	config.MinConns = 3

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, apperrors.ErrPoolCreation
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, apperrors.ErrFailedPing
	}

	fmt.Println("successfully connected to postgres")
	return pool, nil
}