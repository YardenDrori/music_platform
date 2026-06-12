package artists

import (
	"time"

	"github.com/google/uuid"
)

type NewArtistReq struct {
	Name             string     `json:"name"`
	Aliases          []string   `json:"aliases"`
	Description      *string    `json:"description"`
	IsBand           bool       `json:"isBand"`
	LinkToYouTube    *string    `json:"linkToYouTube"`
	LinkToSpotify    *string    `json:"linkToSpotify"`
	LinkToAppleMusic *string    `json:"linkToAppleMusic"`
	OriginDate       *time.Time `json:"originDate"`
	OriginPlace      *string    `json:"originPlace"`
}

type UpdateArtistReq struct {
	ID               uuid.UUID  `json:"-"`
	Name             *string    `json:"name"`
	Description      *string    `json:"description"`
	IsBand           *bool      `json:"isBand"`
	ArtistImageKey   *uuid.UUID `json:"-"`
	ArtistBannerKey  *uuid.UUID `json:"-"`
	LinkToYouTube    *string    `json:"linkToYouTube"`
	LinkToSpotify    *string    `json:"linkToSpotify"`
	LinkToAppleMusic *string    `json:"linkToAppleMusic"`
	OriginDate       *time.Time `json:"birthDate"`
	OriginPlace      *string    `json:"birthPlace"`
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
		OriginDate:       a.OriginDate,
		OriginPlace:      a.OriginPlace,
		AddedAt:          time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        nil,
		Contributions:    []Contribution{},
	}
}
