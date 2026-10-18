package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (handler *UserHandler) HandlerEmailVerification(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DONE1")
	token := r.URL.Query().Get("token")
	email := r.URL.Query().Get("email")

	fmt.Println("token from user: ", token)
	fmt.Println("email from user: ", email)

	fmt.Println("DONE2")
	if token == "" || email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials, all fields are mandatory"})
		return
	}

	fmt.Println("DONE3")
	res, err := handler.service.ServiceEmailVerificationToken(r.Context(), email, token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	if res == true {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"msg": "success"})
	}

}
