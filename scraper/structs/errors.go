package structs

import "errors"

var (
	ErrNotAllowedMethod = errors.New("method not allowed")
	ErrMissingParameter = errors.New("missing parameter")
	ErrNotFound         = errors.New("record not found")
)
