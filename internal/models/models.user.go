package models

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type UserMetadata struct {
	UserID int64
	Data   map[string]any
}

type UserRequest struct {
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type UserRequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        int64          `json:"id"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

type UserLogin struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}
