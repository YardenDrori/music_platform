package artists

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID               uuid.UUID      `json:"id"`
	Name             string         `json:"name"`
	Aliases          []string       `json:"aliases"`
	Description      *string        `json:"description"`
	IsBand           bool           `json:"isBand"`
	ArtistImageUrl   *string        `json:"artistImageUrl"`
	ArtistBannerUrl  *string        `json:"artistBannerUrl"`
	LinkToYouTube    *string        `json:"linkToYouTube"`
	LinkToSpotify    *string        `json:"linkToSpotify"`
	LinkToAppleMusic *string        `json:"linkToAppleMusic"`
	OriginDate       *time.Time     `json:"birthDate"`
	OriginPlace      *string        `json:"birthPlace"`
	AddedAt          time.Time      `json:"addedAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        *time.Time     `json:"-"`
	Contributions    []Contribution `json:"contributions"`
}

type ArtistSummary struct {
	ArtistID                uuid.UUID  `json:"artistId"`
	ArtistName              string     `json:"artistName"`
	ArtistProfilePictureKey *uuid.UUID `json:"artistProfilePictureKey"`
}

type Contribution struct {
	ContributorID         uuid.UUID `json:"contributorId"`
	ContributorName       string    `json:"contributorName"`
	ContributorProfileUrl *string   `json:"contributorProfileUrl"`
	ContributionDate      time.Time `json:"contributionDate"`
}
