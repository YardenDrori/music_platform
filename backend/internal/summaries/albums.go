package summaries

import "github.com/google/uuid"

type AlbumSummary struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	AlbumArtKey *uuid.UUID `json:"albumArtUrl"`
}
