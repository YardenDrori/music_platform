package summaries

import (
	"time"

	"github.com/google/uuid"
)

type SongSummary struct {
	ID          uuid.UUID
	Title       string
	Runtime     time.Duration
	IsPublic    bool
	CoverArtKey *uuid.UUID
}
