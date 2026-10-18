package router

import (
	"encoding/json"
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

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"success": "true"})
		fmt.Println("received request on: ", r.URL.Path)
	})

	accessToken, err := utils.NewAccessTokenClass()
	if err != nil {
		fmt.Println("error in initiating access token class inside router")
	}

	//r.With(middleware.RateLimitMiddleware(redisClient, 3, 5*time.Minute)).Post("/goAuth/register", userHandler.HandlerCreateUser)
	r.Post("/goAuth/register", userHandler.HandlerCreateUser)
	r.With(middleware.RateLimitMiddleware(redisClient, 3, time.Minute)).Post("/goAuth/login", userHandler.HandlerLogin)
	r.With(middleware.AccessTokenVerifyMiddleware(accessToken), middleware.RateLimitMiddleware(redisClient, 3, time.Minute)).Get("/goAuth/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received on health: ", r.URL.Path)
	})

	r.With(middleware.RateLimitMiddleware(redisClient, 3, 5*time.Minute), func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("hello !!!")
			next.ServeHTTP(w, r)
		})
	}).Get("/goAuth/verify-email", userHandler.HandlerEmailVerification)
	return r
}
