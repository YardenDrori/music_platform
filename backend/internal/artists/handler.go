package artists

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"golang.org/x/text/message"

	"github.com/YardenDrori/music-platform/internal/identity"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

// 	NewArtist(ctx context.Context, req NewArtistReq) error
// 	GetArtistsByName(ctx context.Context, name string) ([]*Artist, error)
// 	GetArtistByID(ctx context.Context, id uuid.UUID) (*Artist, error)
// 	UpdateArtistDetails(ctx context.Context, req UpdateArtistReq) error
// 	SoftDeleteArtist(ctx context.Context, id uuid.UUID) error
// 	HardDeleteArtist(ctx context.Context, id uuid.UUID) error


func (h *handler) NewArtist(w http.ResponseWriter, r *http.Request) error {
	requester_id, valid := identity.UserIDFromContext(r.Context())
	if !valid {
		writeError(w, http.StatusUnauthorized, "Unauthenticated")
		return
	}
	var req *NewArtistReq
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.service.NewArtist(r.Context(),*req)
	switch {
		case err == nil:
			break
		case 
	}
}
