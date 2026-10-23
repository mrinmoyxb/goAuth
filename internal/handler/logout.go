package handler

import (
	"encoding/json"
	"net/http"
) 

func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err == nil && cookie.Value != "" {
		err := handler.service.Logout(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "logout failed", http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: "",
		Path: "auth/refresh",
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"msg": "logout successful"})
}


