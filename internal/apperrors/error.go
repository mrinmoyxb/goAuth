package apperrors

type AuthError struct {
	Code       string
	Message    string
	StatusCode int
}

func (e AuthError) Error() string {
	return e.Message
}
