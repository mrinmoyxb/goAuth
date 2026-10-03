package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return rdb, nil;
}