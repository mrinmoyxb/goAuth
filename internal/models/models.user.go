package models

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

type CreateUserInput struct {
	Email    string
	Password string
	Metadata map[string]any
}
