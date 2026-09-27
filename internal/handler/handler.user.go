package handler

import (
	"encoding/json"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"goAuth/internal/services"
	"net/http"
	"errors"
	"log"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, apperrors.ErrInvalidRequest.Message, apperrors.ErrInvalidRequest.StatusCode)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, apperrors.ErrRequiredField.Message, apperrors.ErrRequiredField.StatusCode)
		return
	}

	user, err := handler.service.ServiceCreateUser(r.Context(), req); 
	if err != nil {
		status := http.StatusInternalServerError
		msg := "internal server error"

		switch {
			case errors.Is(err, apperrors.ErrInvalidEmail),
				 errors.Is(err, apperrors.ErrInValidPassword),
				 errors.Is(err, apperrors.ErrEmailAlreadyExists):
				 status = http.StatusBadRequest
				 msg = err.Error()
		}

		http.Error(w, msg, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}