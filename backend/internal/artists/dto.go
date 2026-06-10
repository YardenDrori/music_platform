package artists

import (
	"time"

	"github.com/google/uuid"
)

type NewArtistReq struct {
	Name             string     `json:"name"`
	Description      *string    `json:"description"`
	LinkToYouTube    *string    `json:"linkToYouTube"`
	LinkToSpotify    *string    `json:"linkToSpotify"`
	LinkToAppleMusic *string    `json:"linkToAppleMusic"`
	BirthDate        *time.Time `json:"birthDate"`
	BirthPlace       *string    `json:"birthPlace"`
}

type UpdateArtistReq struct {
	ID               uuid.UUID
	Name             *string    `json:"name"`
	Description      *string    `json:"description"`
	LinkToYouTube    *string    `json:"linkToYouTube"`
	LinkToSpotify    *string    `json:"linkToSpotify"`
	LinkToAppleMusic *string    `json:"linkToAppleMusic"`
	ArtistImageKey   *uuid.UUID `json:"-"`
	ArtistBannerKey  *uuid.UUID `json:"-"`
	BirthDate        *time.Time `json:"birthDate"`
	BirthPlace       *string    `json:"birthPlace"`
	DeletedAt        *time.Time `json:"-"`
}

func (a *NewArtistReq) ToArtist() Artist {
	return Artist{
		ID:               uuid.New(),
		Name:             a.Name,
		Description:      a.Description,
		ArtistImageKey:   nil,
		ArtistBannerKey:  nil,
		LinkToYouTube:    a.LinkToYouTube,
		LinkToSpotify:    a.LinkToSpotify,
		LinkToAppleMusic: a.LinkToAppleMusic,
		BirthDate:        a.BirthDate,
		BirthPlace:       a.BirthPlace,
		AddedAt:          time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        nil,
	}
}
