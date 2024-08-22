package util

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrMalformedRequest    = errors.New("malformed request")
)
