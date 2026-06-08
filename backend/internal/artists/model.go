package artists

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Description      *string    `json:"description"`
	ArtistImageKey   *uuid.UUID `json:"-"`
	ArtistBannerKey  *uuid.UUID `json:"-"`
	LinkToYouTube    *string    `json:"linkToYouTube"`
	LinkToSpotify    *string    `json:"linkToSpotify"`
	LinkToAppleMusic *string    `json:"linkToAppleMusic"`
	BirthDate        *time.Time `json:"birthDate"`
	BirthPlace       *string    `json:"birthPlace"`
	AddedAt          time.Time  `json:"addedAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	DeletedAt        *time.Time `json:"-"`
	UploaderID       uuid.UUID  `json:"uploaderId"`
}
