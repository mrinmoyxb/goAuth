package apperrors

var (
	ErrResetTokenInvalid = AuthError{
		Code:       "invalid_reset_token",
		Message:    "Invalid password reset token",
		StatusCode: 400,
	}

	ErrResetTokenExpired = AuthError{
		Code:       "reset_token_expired",
		Message:    "Password reset token has expired",
		StatusCode: 400,
	}
)
