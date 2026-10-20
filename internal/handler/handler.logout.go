package handler

import "net/http"

func clearRefreshCookie(w http.ResponseWriter){
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: "",
		Path: "/auth/refresh",
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge: -1,
	})
}

func (handler *UserHandler) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err == nil && cookie.Value != "" {

	}
}