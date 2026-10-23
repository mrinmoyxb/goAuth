package handler

import (
	"encoding/json"
	"goAuth/internal/apperrors"
	"net/http"
)

type ResetPassword struct {
	UserID   int64  `json:"userID"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (handler *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	res := &ResetPassword{}

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, apperrors.ErrInvalidCredentials.Message, apperrors.ErrInvalidCredentials.StatusCode)
		return
	}

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
	
	err := handler.service.ResetPassword(r.Context(), res.UserID, res.Token, res.Password)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidCredentials.Message, apperrors.ErrInvalidCredentials.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "token validated"})
}
