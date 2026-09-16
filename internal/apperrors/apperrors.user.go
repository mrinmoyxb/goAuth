package apperrors

var (
	ErrEmailAlreadyExists = AuthError{
		Code: "email_already_exists",
		Message: "An account with this email already exists",
		StatusCode: 409,
	}

	ErrUserNotFound = AuthError{
		Code: "user_not_found",
		Message: "User not found",
		StatusCode: 404,
	}

	ErrInvalidCredentials = AuthError{
		Code: "invalid_credentials",
		Message: "Invalid email or password",
		StatusCode: 401,
	}

	ErrAccountDisabled = AuthError{
		Code: "account_disabled",
		Message: "This account has been disabled",
		StatusCode: 403,
	}

	ErrAccountLocked = AuthError{
		Code: "account_locked",
		Message: "This account is temporarily locker",
		StatusCode: 423,
	}
)