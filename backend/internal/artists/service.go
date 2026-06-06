package artists

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	NewArtist(ctx context.Context, artist Artist) error

	GetArtistsByName(ctx context.Context, name string) ([]Artist, error)
	GetArtistByID(ctx context.Context, id uuid.UUID) (Artist, error)

	UpdateArtist(ctx context.Context, artist NewArtistReq) error
	DeleteArtist(ctx context.Context, id uuid.UUID) error
}
