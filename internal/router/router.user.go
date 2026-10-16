package router

import (
	"fmt"
	"goAuth/internal/handler"
	"goAuth/internal/middleware"
	"goAuth/internal/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func SetUpRouter(userHandler *handler.UserHandler, redisClient *redis.Client) *chi.Mux {
	r := chi.NewRouter()

	accessToken, err := utils.NewAccessTokenClass()
	if err != nil {
		fmt.Println("error in initiating access token class inside router")
	}

	r.With(middleware.RateLimitMiddleware(redisClient, 3, 5*time.Minute)).Post("/goAuth/register", userHandler.HandlerCreateUser)
	r.With(middleware.RateLimitMiddleware(redisClient, 3, time.Minute)).Post("/goAuth/login", userHandler.HandlerLogin)
	r.With(middleware.AccessTokenVerifyMiddleware(accessToken), middleware.RateLimitMiddleware(redisClient, 3, time.Minute)).Get("/goAuth/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received on health: ", r.URL.Path)
	})

	r.With(middleware.RateLimitMiddleware(redisClient, 3, 5*time.Minute)).Post("/goAuth/verify-email", userHandler.HandlerEmailVerification)
	return r
}
