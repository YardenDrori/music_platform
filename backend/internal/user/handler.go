package user

import (
	"fmt"
	"net/http"
)

type handler struct {
	service AuthService
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	panic("")
}
