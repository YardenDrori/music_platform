package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YardenDrori/music-platform/internal/apperrors"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) error {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	resp, newUser, err := h.service.Register(r.Context(), &req)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    resp.refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  resp.refreshExpirey,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(&authResponse{
		User:        newUser,
		AccessToken: resp.accessToken,
	}); err != nil {
		return fmt.Errorf("encoding response", "error", err)
	}
	return nil
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) error {
	var req loginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	resp, user, err := h.service.Login(r.Context(), &req)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    resp.refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  resp.refreshExpirey,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&authResponse{
		User:        user,
		AccessToken: resp.accessToken,
	}); err != nil {
		return fmt.Errorf("encoding response", "error", err)
	}
	return nil
}

func (h *handler) GetAccessToken(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		return apperrors.NewErrBadRequest("refresh-token missing").WithCause(err)
	}

	refreshToken := cookie.Value
	tokens, err := h.service.RequestAccessToken(r.Context(), refreshToken)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    tokens.refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  tokens.refreshExpirey,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&accessTokenResp{
		AccessToken: tokens.accessToken,
	}); err != nil {
		return fmt.Errorf("encoding response", "error", err)
	}
	return nil
}
