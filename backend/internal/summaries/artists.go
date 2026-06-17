package summaries

import "github.com/google/uuid"

type ArtistSummary struct {
	ArtistID                uuid.UUID  `json:"artistId"`
	ArtistName              string     `json:"artistName"`
	ArtistProfilePictureKey *uuid.UUID `json:"artistProfilePictureKey"`
}
