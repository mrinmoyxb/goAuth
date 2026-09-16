package apperrors

type AuthError struct {
	Code       string
	Message    string
	StatusCode int
}

