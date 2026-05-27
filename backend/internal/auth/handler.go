package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func writeInternalError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "internal server error")
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		slog.Error("encoding response", "error", err)
	}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, newUser, err := h.service.Register(r.Context(), &req)
	if err != nil {
		if err, ok := errors.AsType[*ErrBadRequest](err); ok {
			writeError(w, http.StatusBadRequest, err.Message)
			return
		}

		slog.Error("registering new account: ", "error", err)
		writeInternalError(w)
		return
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
		slog.Error("encoding response", "error", err)
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, user, err := h.service.Login(r.Context(), &req)

	if errors.Is(err, ErrUnauthenticated) {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err, ok := errors.AsType[*ErrBadRequest](err); ok {
		writeError(w, http.StatusBadRequest, err.Message)
		return
	}
	if err != nil {
		slog.Error("logging in: ", "error", err)
		writeInternalError(w)
		return
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
		slog.Error("encoding response", "error", err)
	}
}

func (h *handler) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("refresh-token")
	if refreshToken == "" {
		writeError(w, http.StatusBadRequest, "refresh-token header not included")
	}

	tokens, err := h.service.RequestAccessToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, ErrBadToken) {
			writeError(w, http.StatusUnauthorized, "Bad token")
		} else {
			slog.Error("getting access token: ", "error", err)
			writeInternalError(w)
		}
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
		slog.Error("encoding response", "error", err)
	}
}
