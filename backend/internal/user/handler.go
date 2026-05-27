package user

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/YardenDrori/music-platform/internal/identity"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func writeInternalError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "Internal server error")
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		slog.Error("encoding response", "error", err)
	}
}

func (h *handler) GetMe(w http.ResponseWriter, r *http.Request) {
	id, ok := identity.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthenticated")
		return
	}

	user, err := h.service.FindByUUID(r.Context(), id)
	switch {
	case err == nil:
		break
	case errors.Is(err, ErrUnauthenticated):
		writeError(w, http.StatusUnauthorized, "Unauthenticated")
		return
	case errors.Is(err, ErrForbidden):
		writeError(w, http.StatusForbidden, "Forbidden")
		return
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, "User not found")
		return
	default:
		slog.Error("updating user", "error", err)
		writeInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.Error("encoding response", "error", err)
	}
}

func (h *handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	var newAccount NewUserRequest
	if err := json.NewDecoder(r.Body).Decode(&newAccount); err != nil {
		writeError(w, http.StatusBadRequest, "Malformed user provided")
	}

	err := h.service.UpdateAccount(r.Context(), &newAccount)
	switch {
	case errors.Is(err, ErrConflict):
		writeError(w, http.StatusConflict, "Username or Email aren't available")
		return
	case errors.Is(err, ErrForbidden):
		writeError(w, http.StatusForbidden, "Forbidden")
		return
	case errors.Is(err, ErrUnauthenticated):
		writeError(w, http.StatusUnauthorized, "Unauthenticated")
		return
	default:
		if badReq, ok := errors.AsType[*ErrBadRequest](err); ok {
			writeError(w, http.StatusBadRequest, badReq.Message)
			return
		}
		slog.Error("updating user", "error", err)
		writeInternalError(w)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) DisableMe(w http.ResponseWriter, r *http.Request) {
	id, ok := identity.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthenticated")
		return
	}

	err := h.service.DeactivateAccount(r.Context(), id)
	switch {
	case errors.Is(err, ErrUnauthenticated):
		writeError(w, http.StatusUnauthorized, "Unauthenticated")
	case errors.Is(err, ErrForbidden):
		writeError(w, http.StatusForbidden, "Forbidden")
	default:
		writeInternalError(w)
	}

	w.WriteHeader(http.StatusNoContent)
}
