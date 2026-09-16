package apperrors

var (
    ErrInvalidEmail = AuthError{
        Code:       "invalid_email",
        Message:    "Invalid email address",
        StatusCode: 400,
    }

    ErrInvalidPassword = AuthError{
        Code:       "invalid_password",
        Message:    "Invalid password",
        StatusCode: 400,
    }

    ErrPasswordTooShort = AuthError{
        Code:       "password_too_short",
        Message:    "Password is too short",
        StatusCode: 400,
    }

    ErrRequiredField = AuthError{
        Code:       "required_field",
        Message:    "A required field is missing",
        StatusCode: 400,
    }

    ErrInvalidRequest = AuthError{
        Code:       "invalid_request",
        Message:    "Invalid request",
        StatusCode: 400,
    }
)