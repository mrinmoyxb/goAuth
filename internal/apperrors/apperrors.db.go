package apperrors

import "errors"

var (
	ErrDBURLNotFound = errors.New("DATABASE URL NOT FOUND")
	ErrParseDBURL    = errors.New("FAILED TO PARSE DATABASE URL")
	ErrPoolCreation  = errors.New("FAILED TO CREATE POSTGRES POOL")
	ErrFailedPing    = errors.New("FAILED TO PING POSTGRES")
)
