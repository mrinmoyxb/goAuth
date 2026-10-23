package handler

import (
	"encoding/json"
	"goAuth/internal/apperrors"
	"net/http"
)

type RequestPasswordResetToken struct {
	UserID int64 `json:"userID"`
	Token  string `json:"token"`
}

func (handler *UserHandler) ValidatePasswordResetToken(w http.ResponseWriter, r *http.Request) {
	res := &RequestPasswordResetToken{}

	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
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

	err = handler.service.ValidatePasswordResetToken(r.Context(), res.UserID, res.Token)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidCredentials.Message, apperrors.ErrInvalidCredentials.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"msg": "token validated"})
}
