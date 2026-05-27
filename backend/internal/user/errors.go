package user

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

var ErrConflict = errors.New("conflict")

var ErrUnathenticated = errors.New("unauthenticated")

var ErrForbidden = errors.New("forbidden")

type ErrBadRequest struct {
	Message string
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}
