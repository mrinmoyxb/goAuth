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

type EmailVerification struct {
	Email     string
	TokenHash string
	ExpiresAt time.Time
}

type UserMetadata struct {
	UserID int64
	Data   map[string]any
}

type RegisterRequest struct {
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID               int64          `json:"id"`
	Email            string         `json:"email"`
	EmailVerified    bool           `json:"email_verified"`
	EmailStatus      string         `json:"email_status"`
	CreatedAt        time.Time      `json:"created_at"`
}

type LoginResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyEmailRequest struct {
	UserID    int64  `json:"userID"`
	TokenHash string `json:"tokenHash"`
}

type EmailVerificationResponse struct {
	UserID    int64     `json:"userID"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expiresAt"`
}
