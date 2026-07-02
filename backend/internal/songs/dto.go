package songs

import (
	"time"

	"github.com/google/uuid"
)

type NewSongReq struct {
	Title        string      `json:"title"`
	MainArtistID *uuid.UUID  `json:"mainArtistId"`
	AlbumID      *uuid.UUID  `json:"albumId"`
	TrackNumber  *int        `json:"trackNumber"`
	PremieredAt  *time.Time  `json:"premieredAt"`
	SongType     *SongType   `json:"songType"`
	IsPublic     bool        `json:"isPublic"`
	AudioURL     string      `json:"audioURL"`
	ArtistsIDs   []uuid.UUID `json:"artistsIDs"`
}

type UpdateSongReq struct {
	ID           uuid.UUID  `json:"id"`
	Title        *string    `json:"title"`
	MainArtistID *uuid.UUID `json:"mainArtistId"`
	AlbumID      *uuid.UUID `json:"albumId"`
	TrackNumber  *int       `json:"trackNumber"`
	PremieredAt  *time.Time `json:"premieredAt"`
	SongType     *SongType  `json:"songType"`
	IsPublic     *bool      `json:"isPublic"`
	CoverArtKey  *uuid.UUID `json:"coverArtKey"`
	DeletedAt    *time.Time `json:"deletedAt"`

	ArtistsToAdd         []uuid.UUID    `json:"artistsToAdd"`
	ArtistsToRemove      []uuid.UUID    `json:"artistsToRemove"`
	ContributorsToAdd    []Contribution `json:"contributorsToAdd"`
	ContributorsToRemove []uuid.UUID    `json:"contributorsToRemove"`
	WhitelistToAdd       []uuid.UUID    `json:"whitelistToAdd"`
	WhitelistToRemove    []uuid.UUID    `json:"whitelistToRemove"`
}

type GetSongResp struct {
	ID            uuid.UUID      `json:"id"`
	Title         string         `json:"title"`
	MainArtistID  *uuid.UUID     `json:"mainArtistId"`
	AlbumID       *uuid.UUID     `json:"albumId"`
	TrackNumber   *int           `json:"trackNumber"`
	PremieredAt   *time.Time     `json:"premieredAt"`
	Runtime       time.Duration  `json:"runtime"`
	SongType      SongType       `json:"songType"`
	UploadMethod  UploadSource   `json:"uploadMethod"`
	IsPublic      bool           `json:"isPublic"`
	AudioURL      string         `json:"audioUrl"`
	CoverArtURL   *string        `json:"coverArtUrl"`
	AddedAt       time.Time      `json:"addedAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	PlayCount     uint           `json:"playCount"`
	ArtistsIDs    []uuid.UUID    `json:"artistsIds"`
	Contributions []Contribution `json:"contributions"`
}
