package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(rdb *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Invalid client address", http.StatusBadRequest)
				return
			}

			key := "rate_limit:" + host

			count, err := rdb.Incr(r.Context(), key).Result()
			if err != nil {
				http.Error(w, "redis error", http.StatusInternalServerError)
			}

			if count == 1 {
				err := rdb.Expire(r.Context(), key, window).Err()
				if err != nil {
					http.Error(w, "Redis error", http.StatusInternalServerError)
				}
			}

			if count > int64(limit) {
				http.Error(
					w, "Too many requests", http.StatusTooManyRequests,
				)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request recieved middleware")
		next.ServeHTTP(w, r)
	})
}
