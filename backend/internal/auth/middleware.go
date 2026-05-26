package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

type contextClaims struct{}

func NewRequireAuth(
	validator TokenValidator, //takes our validation implementation via validator
) func(http.HandlerFunc) http.HandlerFunc /*returns middleware*/ {

	//we return the middleware
	return func(next http.HandlerFunc) http.HandlerFunc {

		//middleware returns the handler
		return func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok {
				writeError(w, http.StatusBadRequest, "bad token")
				return
			}

			claims, err := validator.ValidateAccessToken(r.Context(), token)
			if errors.Is(err, ErrBadToken) {
				writeError(w, http.StatusUnauthorized, "bad token")
				return
			}
			if msg, ok := errors.AsType[*ErrBadRequest](err); ok {
				writeError(w, http.StatusBadRequest, msg.Message)
				return
			}
			if err != nil {
				slog.Error("failed to validate token", "err", err)
				writeInternalError(w)
				return
			}

			ctx := context.WithValue(r.Context(), contextClaims{}, claims)
			modifiedReq := r.WithContext(ctx)

			next(w, modifiedReq)
		}
	}
}
