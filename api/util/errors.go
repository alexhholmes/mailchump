package util

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidUUID         = errors.New("invalid UUID")
	ErrForbidden           = errors.New("user is not authorized to perform this action")
	ErrMalformedRequest    = errors.New("malformed request")
)
