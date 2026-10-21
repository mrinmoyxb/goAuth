package handler

import (
	"encoding/json"
	"fmt"
	"goAuth/internal/middleware"
	"net/http"
	"strconv"
)

func clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func (handler *UserHandler) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err == nil && cookie.Value != "" {
		err := handler.service.ServiceLogout(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "logout failed", http.StatusInternalServerError)
			return
		}
	}

	clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"msg": "logout successfull"})
}

func (handler *UserHandler) HandlerLogoutAll(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFormContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(user.UserID, 10, 64)
	if err != nil {
		fmt.Println("Error during conversion:", err)
		return
	}

	if err := handler.service.ServiceLogoutAll(r.Context(), userID); err != nil {
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}
}
