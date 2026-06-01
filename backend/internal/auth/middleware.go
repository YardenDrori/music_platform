package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	"github.com/YardenDrori/music-platform/internal/identity"
)

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
			if errors.Is(err, apperrors.ErrBadToken) {
				writeError(w, http.StatusUnauthorized, "bad token")
				return
			}
			if msg, ok := errors.AsType[*apperrors.ErrBadRequest](err); ok {
				writeError(w, http.StatusBadRequest, msg.Message)
				return
			}
			if err != nil {
				slog.Error("failed to validate token", "err", err)
				writeInternalError(w)
				return
			}

			id, err := uuid.Parse(claims.Subject)
			if err != nil {
				slog.Error(
					"claims.Subject is not a valid UUID",
					"subject",
					claims.Subject,
					"err",
					err,
				)
				writeInternalError(w)
				return
			}

			ctx := identity.WithUserID(r.Context(), id)
			modifiedReq := r.WithContext(ctx)

			next(w, modifiedReq)
		}
	}
}
