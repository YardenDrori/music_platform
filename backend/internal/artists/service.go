package artists

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	NewArtist(ctx context.Context, artist Artist, uploaderID uuid.UUID) error

	GetArtistsByNameOrAlias(ctx context.Context, name string) ([]Artist, error)
	GetArtistByID(ctx context.Context, id uuid.UUID) (*Artist, error)

	UpdateArtist(ctx context.Context, req *UpdateArtistReq) error
	DeleteArtist(ctx context.Context, id uuid.UUID) error
}

type Service interface {
	NewArtist(ctx context.Context, req NewArtistReq) error

	GetArtistsByNameOrAlias(ctx context.Context, name string) ([]Artist, error)
	GetArtistByID(ctx context.Context, id uuid.UUID) (*Artist, error)

	UpdateArtistDetails(ctx context.Context, req *UpdateArtistReq) error

	SoftDeleteArtist(ctx context.Context, id uuid.UUID) error
	HardDeleteArtist(ctx context.Context, id uuid.UUID) error

	UploadArtistProfilePicture(ctx context.Context, file []byte, artistID uuid.UUID) error
	UploadArtistBannerPicture(ctx context.Context, file []byte, artistID uuid.UUID) error
}
