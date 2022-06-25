package models

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidJSON = errors.New("invalid json")
	ErrEmptyBody   = errors.New("empty body")
)
