package apperr

import "errors"

var (
	ErrUnauthorized      = errors.New("user is not authorized")
	ErrNotFound          = errors.New("entity not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
