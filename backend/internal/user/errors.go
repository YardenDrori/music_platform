package user

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

type ErrBadRequest struct {
	Message string
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}
