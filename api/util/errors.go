package util

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidUUID         = errors.New("invalid UUID")
	ErrMalformedRequest    = errors.New("malformed request")
)
