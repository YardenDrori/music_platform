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
	ID                   uuid.UUID      `json:"-"`
	Name                 *string        `json:"name"`
	AliasesToAdd         []string       `json:"aliasesToAdd"`
	AliasesToRemove      []string       `json:"aliasesToRemove"`
	Description          *string        `json:"description"`
	IsBand               *bool          `json:"isBand"`
	ArtistImageKey       *uuid.UUID     `json:"-"`
	ArtistBannerKey      *uuid.UUID     `json:"-"`
	LinkToYouTube        *string        `json:"linkToYouTube"`
	LinkToSpotify        *string        `json:"linkToSpotify"`
	LinkToAppleMusic     *string        `json:"linkToAppleMusic"`
	OriginDate           *time.Time     `json:"birthDate"`
	OriginPlace          *string        `json:"birthPlace"`
	DeletedAt            *time.Time     `json:"-"`
	ContributorsToAdd    []Contribution `json:"-"`
	ContributorsToRemove []Contribution `json:"-"`
}

func (a *NewArtistReq) ToArtist() Artist {
	return Artist{
		ID:               uuid.New(),
		Name:             a.Name,
		Description:      a.Description,
		ArtistImageUrl:   nil,
		ArtistBannerUrl:  nil,
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
