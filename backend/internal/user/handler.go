package user

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
