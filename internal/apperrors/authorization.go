package apperrors

var (
	ErrUnauthorized = AuthError{
		Code:       "unauthorized",
		Message:    "Authentication required",
		StatusCode: 401,
	}

	ErrForbidden = AuthError{
		Code:       "forbidden",
		Message:    "You do not have permission to perform this action",
		StatusCode: 403,
	}

	ErrInsufficientPermissions = AuthError{
		Code:       "insufficient_permissions",
		Message:    "Insufficient permissions",
		StatusCode: 403,
	}
)
