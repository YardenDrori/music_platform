package albums

import (
	"time"

	"github.com/google/uuid"
)

type NewAlbumRequest struct {
	Name         string      `json:"name"`
	Description  *string     `json:"description"`
	MainArtistID uuid.UUID   `json:"mainArtistID"`
	HasAllTracks bool        `json:"hasAllTracks"`
	PremieredAt  *time.Time  `json:"premieredAt"`
	Artists      []uuid.UUID `json:"artists"`
}

func (r *NewAlbumRequest) ToAlbum() *Album {
	return &Album{
		ID:           uuid.New(),
		Name:         r.Name,
		Description:  r.Description,
		MainArtistID: r.MainArtistID,
		AlbumArtUrl:  nil,
		HasAllTracks: r.HasAllTracks,
		AddedAt:      time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		DeletedAt:    nil,
		PremieredAt:  r.PremieredAt,
		Artists:      r.Artists,
		Contributors: nil,
	}
}

type UpdateAlbumRequest struct {
	ID                     uuid.UUID   `json:"id"`
	Name                   *string     `json:"name"`
	Description            *string     `json:"description"`
	MainArtistID           *uuid.UUID  `json:"mainArtistId"`
	AlbumArtKey            *uuid.UUID  `json:"albumArtKey"`
	HasAllTracks           *bool       `json:"hasAllTracks"`
	DeletedAt              *time.Time  `json:"-"`
	PremieredAt            *time.Time  `json:"premieredAt"`
	ArtistsToAdd           []uuid.UUID `json:"artistsToAdd"`
	ArtistsToRemove        []uuid.UUID `json:"artistsToRemove"`
	ContributorIDsToAdd    []uuid.UUID `json:"-"`
	ContributorIDsToRemove []uuid.UUID `json:"-"`
}
