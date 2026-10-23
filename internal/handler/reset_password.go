package handler

import (
	"encoding/json"
	"goAuth/internal/apperrors"
	"net/http"
)


func (handler *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	res := &ResetPassword{}

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
