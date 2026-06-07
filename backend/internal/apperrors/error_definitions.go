package apperrors

type errBase struct {
	InternalMessage string
	PublicMessage   string
	Cause           error
}

func (e *errBase) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return e.PublicMessage
}

func (e *errBase) WithInternalMSG(internalErrorMessage string) {
	e.InternalMessage = internalErrorMessage
}

func (e *errBase) WithCause(cause error) {
	e.Cause = cause
}

func (e *errBase) Unwrap() error { return e.Cause }

type ErrUnauthenticated struct{ errBase }

type ErrBadToken struct{ errBase }

type ErrForbidden struct{ errBase }

type ErrNotFound struct{ errBase }

type ErrConflict struct{ errBase }

type ErrBadRequest struct{ errBase }

func NewErrUnathenticated(
	pubErrorMessage string,
) *ErrUnauthenticated {
	return &ErrUnauthenticated{
		errBase{
			PublicMessage: pubErrorMessage,
		},
	}
}

func NewErrBadToken(
	publicMessage string,
) *ErrBadToken {
	return &ErrBadToken{
		errBase{
			PublicMessage: publicMessage,
		},
	}
}

func NewErrForbidden(
	publicMessage string,
) *ErrForbidden {
	return &ErrForbidden{
		errBase{
			PublicMessage: publicMessage,
		},
	}
}

func NewErrNotFound(
	publicMessage string,
) *ErrNotFound {
	return &ErrNotFound{
		errBase{
			PublicMessage: publicMessage,
		},
	}
}

func NewErrConflict(
	publicMessage string,
) *ErrConflict {
	return &ErrConflict{
		errBase{
			PublicMessage: publicMessage,
		},
	}
}

func NewErrBadRequest(
	publicMessage string,
) *ErrBadRequest {
	return &ErrBadRequest{
		errBase{
			PublicMessage: publicMessage,
		},
	}
}
