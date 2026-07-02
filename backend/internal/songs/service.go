package songs

import (
	"context"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/summaries"
)

type Repository interface {
	NewSong(ctx context.Context, song NewSong) error
	GetSongByID(ctx context.Context, id uuid.UUID) (*SongRow, error)
	GetSongSummariesByName(ctx context.Context, songName string, limit int) ([]summaries.SongSummary, error)
	UpdateSong(ctx context.Context, req *UpdateSongReq) (*SongRow, error)
	DeleteSong(ctx context.Context, id uuid.UUID) (*SongRow, error)
}
