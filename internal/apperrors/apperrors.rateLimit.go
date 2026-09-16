package apperrors

var (
    ErrTooManyRequests = AuthError{
        Code:       "too_many_requests",
        Message:    "Too many requests",
        StatusCode: 429,
    }

    ErrTooManyLoginAttempts = AuthError{
        Code:       "too_many_login_attempts",
        Message:    "Too many login attempts",
        StatusCode: 429,
    }
)