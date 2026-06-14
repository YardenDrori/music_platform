package albums

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	NewAlbum(ctx context.Context, album *Album) error
	GetAlbumsByName(ctx context.Context, name string) ([]Album, error)
	GetAlbumByID(ctx context.Context, id uuid.UUID) (*Album, error)
	UpdateAlbum(ctx context.Context, req *UpdateAlbumRequest) (*Album, error)
	DeleteAlbum(ctx context.Context, id uuid.UUID) (*Album, error)
}

type Service interface {
	NewAlbum(ctx context.Context, req *NewAlbumRequest) error

	GetAlbumsByName(ctx context.Context, name string) ([]Album, error)
	GetAlbumByID(ctx context.Context, id uuid.UUID) (*Album, error)

	UpdateAlbumDetails(ctx context.Context, req *UpdateAlbumRequest) error

	SoftDeleteAlbum(ctx context.Context, id uuid.UUID) error
	HardDeleteAlbum(ctx context.Context, id uuid.UUID) error

	UploadAlbumPicture(ctx context.Context, file []byte, artistID uuid.UUID) error
}
