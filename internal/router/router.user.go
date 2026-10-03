package router

import (
	"fmt"
	"goAuth/internal/handler"
	"goAuth/internal/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func SetUpRouter(urlHandler *handler.UserHandler, redisClient *redis.Client) *chi.Mux {
	r := chi.NewRouter()

	r.With(middleware.MiddlewareRateLimit(redisClient, 3, 5*time.Minute)).Post("/goAuth/register", urlHandler.HandlerCreateUser)
	r.With(middleware.MiddlewareRateLimit(redisClient, 3, time.Minute)).Post("/goAuth/login", urlHandler.HandlerLogin)


	r.With(middleware.MiddlewareLogger).Get("/test/health", func(w http.ResponseWriter, r *http.Request){
		fmt.Println("Request received on health: ", r.URL.Path)
	})

	return r
}