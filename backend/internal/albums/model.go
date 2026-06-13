package albums

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	Description  *string        `json:"description"`
	MainArtistID uuid.UUID      `json:"mainArtistId"`
	AlbumArtUrl  *string        `json:"albumArtUrl"`
	HasAllTracks bool           `json:"hasAllTracks"`
	AddedAt      time.Time      `json:"addedAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    *time.Time     `json:"-"`
	PremieredAt  *time.Time     `json:"premieredAt"`
	Artists      []uuid.UUID    `json:"artists"`
	Contributors []Contributor  `json:"contributors"`
}

type Contributor struct {
	ContributorID         uuid.UUID `json:"contributorId"`
	ContributorName       string    `json:"contributorName"`
	ContributorProfileUrl *string   `json:"contributorProfileUrl"`
	ContributionDate      time.Time `json:"contributionDate"`
}
