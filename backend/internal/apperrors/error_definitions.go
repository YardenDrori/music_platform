package apperrors

type errBase struct {
	InternalMessage string
	PublicMessage   string
	Cause           error
}

func (e *errBase) Error() string {
	msg := e.InternalMessage
	if msg == "" {
		msg = e.PublicMessage
	}
	if e.Cause != nil {
		return msg + ": " + e.Cause.Error()
	}
	return msg
}

func (e *errBase) Unwrap() error { return e.Cause }

type ErrUnauthenticated struct{ errBase }

type ErrBadToken struct{ errBase }

type ErrForbidden struct{ errBase }

type ErrNotFound struct{ errBase }

type ErrConflict struct{ errBase }

type ErrBadRequest struct{ errBase }

type ErrInternal struct{ errBase }

func NewErrUnauthenticated(
	pubErrorMessage string,
) *ErrUnauthenticated {
	return &ErrUnauthenticated{
		errBase{
			PublicMessage: pubErrorMessage,
		},
	}
}

func (e *ErrUnauthenticated) WithInternal(msg string) *ErrUnauthenticated {
	e.InternalMessage = msg
	return e
}

func (e *ErrUnauthenticated) WithCause(err error) *ErrUnauthenticated {
	e.Cause = err
	return e
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

func (e *ErrBadToken) WithInternal(msg string) *ErrBadToken {
	e.InternalMessage = msg
	return e
}

func (e *ErrBadToken) WithCause(err error) *ErrBadToken {
	e.Cause = err
	return e
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

func (e *ErrForbidden) WithInternal(msg string) *ErrForbidden {
	e.InternalMessage = msg
	return e
}

func (e *ErrForbidden) WithCause(err error) *ErrForbidden {
	e.Cause = err
	return e
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

func (e *ErrNotFound) WithInternal(msg string) *ErrNotFound {
	e.InternalMessage = msg
	return e
}

func (e *ErrNotFound) WithCause(err error) *ErrNotFound {
	e.Cause = err
	return e
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

func (e *ErrConflict) WithInternal(msg string) *ErrConflict {
	e.InternalMessage = msg
	return e
}

func (e *ErrConflict) WithCause(err error) *ErrConflict {
	e.Cause = err
	return e
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

func (e *ErrBadRequest) WithInternal(msg string) *ErrBadRequest {
	e.InternalMessage = msg
	return e
}

func (e *ErrBadRequest) WithCause(err error) *ErrBadRequest {
	e.Cause = err
	return e
}

func NewErrInternal() *ErrInternal {
	return &ErrInternal{
		errBase{
			PublicMessage: "internal server error",
		},
	}
}

func (e *ErrInternal) WithInternal(msg string) *ErrInternal {
	e.InternalMessage = msg
	return e
}

func (e *ErrInternal) WithCause(err error) *ErrInternal {
	e.Cause = err
	return e
}
