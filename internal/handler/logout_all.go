package handler

import (
	"fmt"
	"goAuth/internal/middleware"
	"net/http"
	"strconv"
) 

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