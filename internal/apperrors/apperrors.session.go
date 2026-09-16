package apperrors

var (
    ErrSessionNotFound = AuthError{
        Code:       "session_not_found",
        Message:    "Session not found",
        StatusCode: 401,
    }

    ErrSessionExpired = AuthError{
        Code:       "session_expired",
        Message:    "Session has expired",
        StatusCode: 401,
    }

    ErrSessionRevoked = AuthError{
        Code:       "session_revoked",
        Message:    "Session has been revoked",
        StatusCode: 401,
    }

    ErrRefreshTokenInvalid = AuthError{
        Code:       "invalid_refresh_token",
        Message:    "Invalid refresh token",
        StatusCode: 401,
    }
)