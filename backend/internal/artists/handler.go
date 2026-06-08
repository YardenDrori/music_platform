package artists

import (
	"encoding/json"
	"net/http"

	"github.com/YardenDrori/music-platform/internal/apperrors"
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
	req := &NewArtistReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperrors.NewErrBadRequest("invalid request body")
	}

	return h.service.NewArtist(r.Context(), *req)
}
