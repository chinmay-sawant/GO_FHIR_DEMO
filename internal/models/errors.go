package models

import "errors"

var (
	// ErrNotFound indicates that the requested resource does not exist.
	ErrNotFound = errors.New("resource not found")
	// ErrInternal indicates an unexpected internal failure.
	ErrInternal = errors.New("internal server error")
	// ErrInvalid indicates that the supplied input is invalid.
	ErrInvalid = errors.New("invalid input")
)
