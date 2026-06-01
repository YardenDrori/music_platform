package apperrors

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

var ErrConflict = errors.New("conflict")

var ErrUnauthenticated = errors.New("incorrect username or password")

var ErrBadToken = errors.New("token invalid or expired")

var ErrForbidden = errors.New("forbidden")

type ErrBadRequest struct {
	Message string
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}
