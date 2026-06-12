package albums

import (
	"time"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/artists"
)

type Album struct {
	Id           uuid.UUID
	Name         string
	Description  *string
	MainArtistID uuid.UUID
	AlbumArtKey  *uuid.UUID
	HasAllTracks bool
	AddedAt      time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	PremieredAt  *time.Time
	Artists      []artists.Artist
	Contributors []Contributor
}

type Contributor struct {
	ContributorID         uuid.UUID `json:"contributorId"`
	ContributorName       string    `json:"contributorName"`
	ContributorProfileUrl *string   `json:"contributorProfileUrl"`
	ContributionDate      time.Time `json:"contributionDate"`
}
