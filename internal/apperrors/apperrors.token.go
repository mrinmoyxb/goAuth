package apperrors

var (
	ErrInvalidToken = AuthError{
		Code:       "invalid_token",
		Message:    "Invalid authentication token",
		StatusCode: 401,
	}

	ErrExpiredToken = AuthError{
		Code:       "token_expired",
		Message:    "Authentication token has expired",
		StatusCode: 401,
	}

	ErrMalformedToken = AuthError{
		Code:       "malformed_token",
		Message:    "Malformed authentication token",
		StatusCode: 401,
	}

	ErrTokenRevoked = AuthError{
		Code:       "token_revoked",
		Message:    "Authentication token has been revoked",
		StatusCode: 401,
	}
)
