package handler

import (
	"encoding/json"
	"goAuth/internal/apperrors"
	"net/http"
	"fmt"
)

type RequestPasswordReset struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	UserID   int64
	Token    string
	Password string
}

func (handler *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) { 
	req := RequestPasswordReset{}

	fmt.Println("EMAIL: ", req.Email)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidRequest.Message, apperrors.ErrInvalidRequest.StatusCode)
		return
	}
	if req.Email == "" {
		http.Error(w, apperrors.ErrInvalidCredentials.Message, apperrors.ErrInvalidCredentials.StatusCode)
		return
	}

	fmt.Println("Check1")
	user, err := handler.service.RequestPasswordReset(r.Context(), req.Email)
	if err != nil {
		http.Error(w, apperrors.ErrInvalidRequest.Message, apperrors.ErrInvalidRequest.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"msg": "sent a token to the mail", "user":user.UserID})
}



