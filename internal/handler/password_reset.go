package handler

import (
	"encoding/json"
	"goAuth/internal/apperrors"
	"net/http"
)

type PasswordResetTokenRequest struct {
	Email string
	Token string
}

type ResetPasswordRequest struct {
	UserID   int64
	Token    string
	Password string
}

func (handler *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var email string

	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidRequest.Message, apperrors.ErrInvalidRequest.StatusCode)
		return
	}
	if email == "" {
		http.Error(w, apperrors.ErrInvalidCredentials.Message, apperrors.ErrInvalidCredentials.StatusCode)
		return
	}

	err = handler.service.RequestPasswordReset(r.Context(), email)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidRequest.Message, apperrors.ErrInvalidRequest.StatusCode)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"msg": "sent a token to the mail"})
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) ValidatePasswordResetToken(w http.ResponseWriter, r *http.Request) {
	res := &PasswordResetTokenRequest{}

	// 1. validate email and token
	err := json.NewDecoder(r.Body).Decode(res)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidCredentials.Message, apperrors.ErrInvalidCredentials.StatusCode)
		return
	}
	if res.Email == "" {
		http.Error(w, apperrors.ErrInvalidEmail.Message, apperrors.ErrInvalidEmail.StatusCode)
		return
	}
	if res.Token == "" {
		http.Error(w, apperrors.ErrInvalidToken.Message, apperrors.ErrInvalidToken.StatusCode)
		return
	}

	// 2. validate token
	user, err := handler.service.ValidatePasswordResetToken(r.Context(), res.Email, res.Token)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidCredentials.Message, apperrors.ErrInvalidCredentials.StatusCode)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"msg": "token validated", "userID": user.UserID})
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	res := &ResetPasswordRequest{}

	// 1. validate email, token and new password
	if res.UserID == 0 {
		http.Error(w, apperrors.ErrInvalidUserID.Message, apperrors.ErrInvalidUserID.StatusCode)
		return
	}
	if res.Token == "" {
		http.Error(w, apperrors.ErrInvalidToken.Message, apperrors.ErrInvalidToken.StatusCode)
		return
	}
	if res.Password == "" {
		http.Error(w, apperrors.ErrInvalidPassword.Message, apperrors.ErrAccountLocked.StatusCode)
		return
	}

	// 2. validate token
	err := handler.service.ResetPassword(r.Context(), res.UserID, res.Token, res.Password)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidCredentials.Message, apperrors.ErrInvalidCredentials.StatusCode)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"msg": "token validated"})
	w.WriteHeader(http.StatusOK)
}
