package artists

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID               uuid.UUID      `json:"id"`
	Name             string         `json:"name"`
	Aliases          []string       `json:"aliases"`
	Description      *string        `json:"description"`
	IsBand           bool           `json:"isBand"`
	ArtistImageKey   *uuid.UUID     `json:"-"`
	ArtistBannerKey  *uuid.UUID     `json:"-"`
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

type ArtistEntity struct {
	Artist
	ContributionsEntity []ContributionEntity
}

// db entity
type ContributionEntity struct {
	ContributorID         uuid.UUID
	ContributorName       string
	ContributorProfileKey *uuid.UUID
	ContributionDate      time.Time
}

// model+dto
type Contribution struct {
	ContributorID         uuid.UUID `json:"contributorId"`
	ContributorName       string    `json:"contributorName"`
	ContributorProfileUrl *string   `json:"contributorProfileUrl"`
	ContributionDate      time.Time `json:"contributionDate"`
}

func (c *ContributionEntity) ToModel(profilePicUrl *url.URL) Contribution {
	var picUrl *string
	if profilePicUrl != nil {
		s := profilePicUrl.String()
		picUrl = &s
	}
	return Contribution{
		ContributorID:         c.ContributorID,
		ContributorName:       c.ContributorName,
		ContributorProfileUrl: picUrl,
		ContributionDate:      c.ContributionDate,
	}
}
