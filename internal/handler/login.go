package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"log"
	"net/http"
)

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	fmt.Println("Check1")

	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, apperrors.ErrInvalidRequest.Message, apperrors.ErrInvalidRequest.StatusCode)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, apperrors.ErrRequiredField.Message, apperrors.ErrRequiredField.StatusCode)
		return
	}

	fmt.Println("Check2")
	user, err := handler.service.Login(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "internal server error"

		switch {
		case errors.Is(err, apperrors.ErrInvalidEmail),
			errors.Is(err, apperrors.ErrInvalidPassword):
			status = http.StatusBadRequest
			msg = err.Error()
		case errors.Is(err, apperrors.ErrInvalidCredentials):
			status = http.StatusBadRequest
			msg = err.Error()
		}

		http.Error(w, msg, status)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    user.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/auth/refresh",
		MaxAge:   15 * 24 * 60 * 60,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"access_token": user.AccessToken, "msg": "login successfull"}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
