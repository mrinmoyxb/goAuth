package apperrors

var (
	ErrEmailNotVerified = AuthError{
		Code:       "email_not_verified",
		Message:    "Email address has not been verified",
		StatusCode: 403,
	}

	ErrVerificationTokenInvalid = AuthError{
		Code:       "invalid_verification_token",
		Message:    "Invalid verification token",
		StatusCode: 400,
	}

	ErrVerificationTokenExpired = AuthError{
		Code:       "verification_token_expired",
		Message:    "Verification token has expired",
		StatusCode: 400,
	}
)
