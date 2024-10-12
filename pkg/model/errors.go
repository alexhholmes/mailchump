package model

import "errors"

var (
	ErrAlreadyExists = errors.New("row already exists")
	ErrNotFound      = errors.New("row not found")
	ErrNoChanges     = errors.New("no changes made")
	ErrTxBegin       = errors.New("failed to begin transaction")
	ErrTxCommit      = errors.New("failed to commit transaction")
)
