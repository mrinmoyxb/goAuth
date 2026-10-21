package apperrors

import "errors"

var (
	ErrFailedToHashPass    = errors.New("FAILED TO HASH PASSWORD")
	ErrInvalidHashFormat   = errors.New("INVALID HASH FORMAT")
	ErrIncompatibleVersion = errors.New("INCOMPATIBLE ARGON2 VERSION")
	ErrInValidPassword     = errors.New("INVALID PASSWORD")
)
