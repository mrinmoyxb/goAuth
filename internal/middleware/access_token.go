package middleware

import (
	"context"
	"goAuth/internal/utils"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "authenticatedUser"

func RequireAccessToken(tokenService *utils.TokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, tokenErr := tokenService.VerifyAccessToken(parts[1])
			if tokenErr != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, *user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) (utils.AuthenticatedUser, bool) {
	user, ok := ctx.Value(userContextKey).(utils.AuthenticatedUser)
	return user, ok
}
