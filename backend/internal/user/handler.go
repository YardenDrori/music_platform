package user

import (
	"encoding/json"
	"fmt"
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
	id, err := identity.UserIDFromContext(r.Context())
	if err != nil {
		return fmt.Errorf("getting own account details: %w", err)
	}

	user, err := h.service.FindByUUIDInternal(r.Context(), id)
	if err != nil {
		e := apperrors.NewErrInternal("").
			WithInternal("failed to get account details despite valid token").
			WithCause(err)
		return fmt.Errorf("getting own account details: %w", e)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return fmt.Errorf("encoding response: %w", err)
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
	id, err := identity.UserIDFromContext(r.Context())
	if err != nil {
		return fmt.Errorf("disabling own account: %w", err)
	}

	err = h.service.DeactivateAccount(r.Context(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
