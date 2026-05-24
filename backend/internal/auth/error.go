package auth

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

var ErrConflict = errors.New("conflict")

var ErrUnauthorized = errors.New("incorrect username or password")

var ErrBadToken = errors.New("token invalid or expired")

type ErrBadRequest struct {
	Message string
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}
