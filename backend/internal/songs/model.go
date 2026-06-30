package songs

import (
	"time"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/summaries"
)

type UploadSource string

type SongType string

const (
	UploadSourceManual     UploadSource = "manual_upload"
	UploadSourceYoutube    UploadSource = "youtube_scrape"
	UploadSourceSpotify    UploadSource = "spotify_scrape"
	UploadSourceAppleMusic UploadSource = "apple_music_scrape"
	UploadSourceOther      UploadSource = "other"
)

const (
	SongTypeAlbumTrack SongType = "album_track"
	SongTypeSingle     SongType = "single"
	SongTypeOrphaned   SongType = "orphaned"
)

type SongRow struct {
	ID           uuid.UUID
	Title        string
	MainArtistID *uuid.UUID
	AlbumID      *uuid.UUID
	TrackNumber  *int
	PremieredAt  *time.Time
	Runtime      time.Duration
	SongType     SongType
	UploadMethod UploadSource
	IsPublic     bool
	AudioKey     uuid.UUID
	CoverArtKey  *uuid.UUID
	AddedAt      time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type Contribution struct {
	UserSummary      summaries.UserSummary
	ContributionDate time.Time
}

// all relevant data for this type
type GetSongDetails struct {
	SongRow
	PlayCount        int
	Album            *summaries.AlbumSummary
	Artists          []summaries.ArtistSummary
	Contributions    []Contribution
	WhitelistedUsers []summaries.UserSummary
}

type NewSong struct {
	SongRow
	Album              *summaries.AlbumSummary
	ArtistIDs          []uuid.UUID
	Contributions      []Contribution
	WhitelistedUserIDs []uuid.UUID
}
