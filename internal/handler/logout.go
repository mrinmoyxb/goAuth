package handler

import (
	"encoding/json"
	"fmt"
	"goAuth/internal/middleware"
	"net/http"
	"strconv"
)

func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err == nil && cookie.Value != "" {
		err := handler.service.Logout(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "logout failed", http.StatusInternalServerError)
			return
		}
	}

	//
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"msg": "logout successfull"})
}

func (handler *UserHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(user.UserID, 10, 64)
	if err != nil {
		fmt.Println("Error during conversion:", err)
		return
	}

	if err := handler.service.LogoutAll(r.Context(), userID); err != nil {
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}
}
