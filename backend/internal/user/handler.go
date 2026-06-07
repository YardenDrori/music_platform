package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	"github.com/YardenDrori/music-platform/internal/identity"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) GetMe(w http.ResponseWriter, r *http.Request) error {
	id, ok := identity.UserIDFromContext(r.Context())
	if !ok {
		return apperrors.NewErrUnathenticated("", "Unauthenticated", nil)
	}

	user, err := h.service.FindByUUID(r.Context(), id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return fmt.Errorf("encoding response", "error", err)
	}
	return nil
}

func (h *handler) UpdateMe(w http.ResponseWriter, r *http.Request) error {
	var newAccount NewUserRequest
	if err := json.NewDecoder(r.Body).Decode(&newAccount); err != nil {
		return apperrors.NewErrBadRequest("malformed user provided")
	}

	err := h.service.UpdateAccount(r.Context(), &newAccount)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *handler) DisableMe(w http.ResponseWriter, r *http.Request) error {
	id, ok := identity.UserIDFromContext(r.Context())
	if !ok {
		return apperrors.NewErrUnathenticated("unauthenticated")
	}

	err := h.service.DeactivateAccount(r.Context(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
