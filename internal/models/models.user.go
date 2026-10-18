package models

import "time"

type User struct {
	ID            int64
	Email         string
	PasswordHash  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	EmailVerified bool
}

// for creation of user
type UserInput struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	TokenHash    string
}

type UserEmailVerify struct {
	Email     string
	TokenHash string
	ExpiresAt time.Time
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
	ID               int64          `json:"id"`
	Email            string         `json:"email"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	Metadata         map[string]any `json:"metadata,omitempty"`
	AccessToken      string         `json:"access_token"`
	VerificationMail bool           `json:"verification_mail"`
}

type LoginResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserEmailVerification struct {
	UserID    int64  `json:"userID"`
	TokenHash string `json:"tokenHash"`
}

type UserEmailVerificationResponse struct {
	UserID    int64     `json:"userID"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expiresAt"`
}
