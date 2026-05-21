package user

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type handler struct {
	service AuthService
}

func NewHandler(service AuthService) *handler {
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

	resp, err := h.service.Register(r.Context(), &req)
	if err != nil {
		if err, ok := errors.AsType[*ErrBadRequest](err); ok {
			writeError(w, http.StatusBadRequest, err.Message)
			return
		}

		slog.Error("registering new account: ", "error", err)
		writeInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
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

	resp, err := h.service.Login(r.Context(), &req)

	if errors.Is(err, ErrUnauthorized) {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		slog.Error("encoding response", "error", err)
	}
}
