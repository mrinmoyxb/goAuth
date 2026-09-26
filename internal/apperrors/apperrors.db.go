package apperrors

import "errors"

var (
	ErrDBURLNotFound = errors.New("DATABASE URL NOT FOUND")
	ErrParseDBURL    = errors.New("FAILED TO PARSE DATABASE URL")
	ErrPoolCreation  = errors.New("FAILED TO CREATE POSTGRES POOL")
	ErrFailedPing    = errors.New("FAILED TO PING POSTGRES")
	ErrTxBeginFailed = errors.New("FAILED TO START TRANSACTION")
	ErrInsertFailed  = errors.New("FAILED TO INSERT DATA")
	ErrQueryFailed   = errors.New("FAILED TO EXECUTE QUERY")
)
