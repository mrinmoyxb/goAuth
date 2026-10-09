package handler

import (
	"encoding/json"
	"errors"
	"goAuth/internal/apperrors"
	"goAuth/internal/models"
	"goAuth/internal/services"
	"log"
	"net/http"
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
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req models.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, apperrors.ErrInvalidRequest.Message, apperrors.ErrInvalidRequest.StatusCode)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, apperrors.ErrRequiredField.Message, apperrors.ErrRequiredField.StatusCode)
		return
	}

	user, err := handler.service.ServiceCreateUser(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "internal server error"

		switch {
		case errors.Is(err, apperrors.ErrInvalidEmail):
			status, msg = http.StatusBadRequest, "invalid email"
		case errors.Is(err, apperrors.ErrInvalidPassword):
			status, msg = http.StatusBadRequest, "invalid password"
		case errors.Is(err, apperrors.ErrEmailAlreadyExists):
			status, msg = http.StatusConflict, "email already exists"
		default:
			log.Printf("register failed: %v", err)
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

func (handler *UserHandler) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req models.UserRequestLogin

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, apperrors.ErrInvalidRequest.Message, apperrors.ErrInvalidRequest.StatusCode)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, apperrors.ErrRequiredField.Message, apperrors.ErrRequiredField.StatusCode)
		return
	}

	user, err := handler.service.ServiceLogin(r.Context(), req)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
