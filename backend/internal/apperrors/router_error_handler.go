package apperrors

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func writeInternalError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "Internal server error")
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}

func HandlerError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		slog.Error("HandlerError called with nil error", "method", r.Method, "path", r.URL.Path)
		writeInternalError(w)
		return
	}
	//this MUST be the first error we check cause an internal server error could
	//have other errors wrapped inside of it
	if e, ok := errors.AsType[*ErrInternal](err); ok {
		slog.Info("internal", "method", r.Method, "path", r.URL.Path, "error", e)
		resolveError(w, http.StatusInternalServerError, "internal error", &e.errBase)
		return
	}
	if e, ok := errors.AsType[*ErrUnauthenticated](err); ok {
		slog.Info("unauthenticated", "method", r.Method, "path", r.URL.Path, "error", e)
		resolveError(w, http.StatusUnauthorized, "unauthenticated", &e.errBase)
		return
	}
	if e, ok := errors.AsType[*ErrBadToken](err); ok {
		slog.Info("bad token", "method", r.Method, "path", r.URL.Path, "error", e)
		resolveError(w, http.StatusUnauthorized, "bad token", &e.errBase)
		return
	}
	if e, ok := errors.AsType[*ErrForbidden](err); ok {
		slog.Info("forbidden", "method", r.Method, "path", r.URL.Path, "error", e)
		resolveError(w, http.StatusForbidden, "forbidden", &e.errBase)
		return
	}
	if e, ok := errors.AsType[*ErrNotFound](err); ok {
		resolveError(w, http.StatusNotFound, "not found", &e.errBase)
		return
	}
	if e, ok := errors.AsType[*ErrConflict](err); ok {
		if e.InternalMessage != "" {
			slog.Info("conflict", "method", r.Method, "path", r.URL.Path, "error", e)
		}
		resolveError(w, http.StatusConflict, "conflict", &e.errBase)
		return
	}
	if e, ok := errors.AsType[*ErrBadRequest](err); ok {
		resolveError(w, http.StatusBadRequest, "bad request", &e.errBase)
		return
	}
	slog.Error("unhandled error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeInternalError(w)
}

func resolveError(
	w http.ResponseWriter,
	status int,
	fallbackMessage string,
	err *errBase,
) {
	message := err.PublicMessage
	if message == "" {
		message = fallbackMessage
	}
	writeError(w, status, message)
}
