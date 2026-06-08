package auth

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	"github.com/YardenDrori/music-platform/internal/identity"
)

func NewRequireAuth(
	validator TokenValidator,
) func(http.ResponseWriter, *http.Request) (http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request, error) {
		header := r.Header.Get("Authorization")
		if header == "" {
			return w, r, apperrors.NewErrUnauthenticated("Authorization header required")
		}

		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok {
			return w, r, apperrors.NewErrBadRequest("Authorization header must use Bearer scheme")
		}

		claims, err := validator.ValidateAccessToken(r.Context(), token)
		if err != nil {
			return w, r, err
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			e := apperrors.NewErrBadToken("bad token").
				WithInternal("user tried to login with a valid token that has invalid uuid syntax").
				WithCause(err)
			return w, r, e
		}

		newContext := identity.WithUserID(r.Context(), userID)
		newReq := r.WithContext(newContext)
		return w, newReq, nil
	}
}
