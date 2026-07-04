package songs

import (
	"context"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/summaries"
)

type repository interface {
	NewSong(ctx context.Context, song NewSong) error
	GetSongByID(ctx context.Context, id uuid.UUID) (*SongRow, error)
	GetSongSummariesByName(
		ctx context.Context,
		songName string,
		limit int,
	) ([]summaries.SongSummary, error)
	UpdateSong(ctx context.Context, req *UpdateSongReq) (*SongRow, error)
	DeleteSong(ctx context.Context, id uuid.UUID) (*SongRow, error)
	NewStagingSong(ctx context.Context, objectKey string, uploadID string, userID uuid.UUID) error
	StagingSongOwnerFromObjectKey(ctx context.Context, uploadID string) (uuid.UUID, error)
	DeleteStagingSong(ctx context.Context, objectKey string) error
	EvictStagingSongs(
		ctx context.Context,
		abort func(objectKey string, uploadID string) error,
	) error
}

type Service interface {
}
