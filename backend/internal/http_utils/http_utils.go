package httputils

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return err
	}
	return nil
}
