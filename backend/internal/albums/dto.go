package albums

import (
	"time"

	"github.com/google/uuid"
)

type UpdateAlbumRequest struct {
	ID                     uuid.UUID
	Name                   *string
	Description            *string
	MainArtistID           *uuid.UUID
	AlbumArtKey            *uuid.UUID
	HasAllTracks           *bool
	DeletedAt              *time.Time
	PremieredAt            *time.Time
	ArtistsToAdd           []uuid.UUID
	ArtistsToRemove        []uuid.UUID
	ContributorIDsToAdd    []uuid.UUID
	ContributorIDsToRemove []uuid.UUID
}
