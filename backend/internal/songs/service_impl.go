package songs

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	"github.com/YardenDrori/music-platform/internal/constants"
	"github.com/YardenDrori/music-platform/internal/identity"
	"github.com/YardenDrori/music-platform/internal/storage"
)

type service struct {
	repo    repository
	storage storage.Service
}

func NewService(repo repository, storage storage.Service) Service {
	return &service{repo: repo, storage: storage}
}

func (s *service) NewSongInit(ctx context.Context) (string, uuid.UUID, error) {
	userID, err := identity.UserIDFromContext(ctx)
	if err != nil {
		return "", uuid.Nil, fmt.Errorf("initiating new song request: %w", err)
	}

	objectKey := uuid.New()
	uploadID, err := s.storage.InitiateMultipartUpload(
		ctx,
		constants.SongDataStagingBucket,
		objectKey.String(),
		storage.PutOptions{
			ContentType: "application/octet-stream",
			Checksum:    storage.ChecksumSha256,
		},
	)
	if err != nil {
		return "", uuid.Nil, fmt.Errorf("initiating new song request: %w", err)
	}

	if err := s.repo.NewStagingSong(ctx, objectKey.String(), uploadID, userID); err != nil {
		if anotherErr := s.storage.AbortMultipartUpload(
			ctx,
			constants.SongDataStagingBucket,
			objectKey.String(),
			uploadID,
		); anotherErr != nil {
			return "", uuid.Nil, fmt.Errorf("initiating new song request: %w, %w", err, anotherErr)
		}
		return "", uuid.Nil, fmt.Errorf("initiating new song request: %w", err)
	}

	return uploadID, objectKey, nil
}

func (s *service) verifyIdentity(ctx context.Context, objectKey string) error {
	userID, err := identity.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("handling new song request: %w", err)
	}
	uploaderID, err := s.repo.StagingSongOwnerFromObjectKey(ctx, objectKey)
	if err != nil {
		return fmt.Errorf("handling new song request: %w", err)
	}
	if uploaderID != userID {
		return fmt.Errorf(
			"handling new song request: %w",
			apperrors.NewErrForbidden("forbidden"),
		)
	}
	return nil
}

func (s *service) NewSongGetPresignedURLs(
	ctx context.Context,
	objectKey string,
	uploadID string,
	totalPartsCount int,
	checksums ...string,
) ([]string, error) {
	if err := s.verifyIdentity(ctx, objectKey); err != nil {
		return nil, err
	}

	if res, err := s.storage.PresignMultipartUploadPutURLs(
		ctx,
		constants.SongDataStagingBucket,
		objectKey,
		uploadID,
		totalPartsCount,
		checksums...); err != nil {
		return nil, fmt.Errorf("handling new song request: %w", err)
	} else {
		return res, nil
	}
}

func (s *service) NewSongCompleteUpload(
	ctx context.Context,
	objectKey string,
	uploadID string,
	completedPartsDTOs []storage.CompletedPartDTO,
) error {
	if err := s.verifyIdentity(ctx, objectKey); err != nil {
		return err
	}

	completedParts := []storage.CompletedPart{}
	for _, part := range completedPartsDTOs {
		newPart := storage.CompletedPart{
			PartNumber:    part.PartNumber,
			ETag:          part.ETag,
			ChecksumValue: part.ChecksumValue,
			ChecksumAlgo:  storage.ChecksumSha256,
		}
		completedParts = append(completedParts, newPart)
	}

	if err := s.storage.CompleteMultipartUpload(
		ctx,
		constants.SongDataStagingBucket,
		objectKey,
		uploadID,
		completedParts,
		storage.PutOptions{}, //options already provided on init
	); err != nil {
		return fmt.Errorf("handling new song request: %w", err)
	}

	if err := s.finalizeStagedSong(ctx, objectKey); err != nil {
		return fmt.Errorf("handling new song request: %w", err)
	}

	if err := s.repo.DeleteStagingSong(ctx, objectKey); err != nil {
		slog.Error("handling new song request", "error", err)
	}

	return nil
}

// todo
func (s *service) finalizeStagedSong(ctx context.Context, objectKey string) error {
	panic("")
}
