package apperrors

import "errors"

var (
	ErrFailedToHashPass = errors.New("FAILED TO HASH PASSWORD")
)