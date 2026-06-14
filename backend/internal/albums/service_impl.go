package albums

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	"github.com/YardenDrori/music-platform/internal/constants"
	"github.com/YardenDrori/music-platform/internal/identity"
	"github.com/YardenDrori/music-platform/internal/storage"
)

type service struct {
	repo    Repository
	storage storage.Service
}

func NewService(repo Repository, storage storage.Service) Service {
	return &service{repo: repo, storage: storage}
}

var allowedImageFormats = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
}

func (s *service) processAlbumEntity(album *Album) {
	if album.AlbumArtUrl != nil {
		album.AlbumArtUrl = s.storage.BuildPublicGetUrl(
			constants.AlbumArtBucket,
			*album.AlbumArtUrl,
		)
	}
	for i := range album.Contributors {
		c := &album.Contributors[i]
		if c.ContributorProfileUrl != nil {
			c.ContributorProfileUrl = s.storage.BuildPublicGetUrl(
				constants.ProfilePicBucket,
				*c.ContributorProfileUrl,
			)
		}
	}
}

func (s *service) validateNewAlbumReq(req *NewAlbumRequest) error {
	if len(req.Artists) == 0 {
		return apperrors.NewErrBadRequest("album must have at least one artist")
	}
	if req.Name == "" {
		return apperrors.NewErrBadRequest("album name cannot be empty")
	}
	if req.PremieredAt != nil && req.PremieredAt.After(time.Now()) {
		return apperrors.NewErrBadRequest("Premiered At date cannot be in the future")
	}
	return nil
}

func validateUpdateAlbumReq(req *UpdateAlbumRequest) error {
	if req.Name != nil && *req.Name == "" {
		return apperrors.NewErrBadRequest("album name cannot be empty")
	}
	if req.PremieredAt != nil && req.PremieredAt.After(time.Now()) {
		return apperrors.NewErrBadRequest("premiered date cannot be in the future")
	}
	return nil
}

func (s *service) NewAlbum(ctx context.Context, req *NewAlbumRequest) error {
	requesterID, err := identity.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("crating new album: %w", err)
	}

	if err := s.validateNewAlbumReq(req); err != nil {
		return fmt.Errorf("creating new album: %w", err)
	}

	request := req.ToAlbum()
	request.Contributors = append(request.Contributors, Contributor{
		ContributorID: requesterID,
	})

	if err := s.repo.NewAlbum(ctx, request); err != nil {
		return fmt.Errorf("creating new album: %w", err)
	}
	return nil
}

func (s *service) GetAlbumsByName(ctx context.Context, name string) ([]Album, error) {
	albums, err := s.repo.GetAlbumsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("getting albums by name: %w", err)
	}
	for i := range albums {
		s.processAlbumEntity(&albums[i])
	}
	return albums, nil
}

func (s *service) GetAlbumByID(ctx context.Context, id uuid.UUID) (*Album, error) {
	album, err := s.repo.GetAlbumByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting album by id: %w", err)
	}
	s.processAlbumEntity(album)
	return album, nil
}

func (s *service) UpdateAlbumDetails(ctx context.Context, req *UpdateAlbumRequest) error {
	requesterID, err := identity.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("updating album details: %w", err)
	}
	if err := validateUpdateAlbumReq(req); err != nil {
		return fmt.Errorf("updating album details: %w", err)
	}
	req.ContributorIDsToAdd = append(req.ContributorIDsToAdd, requesterID)
	if err := s.repo.UpdateAlbum(ctx, req); err != nil {
		return fmt.Errorf("updating album details: %w", err)
	}
	return nil
}

func (s *service) SoftDeleteAlbum(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	if err := s.repo.UpdateAlbum(ctx, &UpdateAlbumRequest{ID: id, DeletedAt: &now}); err != nil {
		return fmt.Errorf("soft deleting album: %w", err)
	}
	return nil
}

func (s *service) HardDeleteAlbum(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteAlbum(ctx, id); err != nil {
		return fmt.Errorf("hard deleting album: %w", err)
	}
	return nil
}

func (s *service) UploadAlbumPicture(ctx context.Context, file []byte, albumID uuid.UUID) error {
	contentType := http.DetectContentType(file)
	if _, ok := allowedImageFormats[contentType]; !ok {
		return fmt.Errorf(
			"uploading album art: %w",
			apperrors.NewErrBadRequest("invalid content type"),
		)
	}

	reader := bytes.NewReader(file)
	config, _, err := image.DecodeConfig(reader)
	if err != nil {
		return fmt.Errorf(
			"uploading album art: %w",
			apperrors.NewErrBadRequest("could not decode image"),
		)
	}
	reader.Reset(file)

	if config.Width != config.Height || config.Width < 512 || config.Width > 4096 {
		return fmt.Errorf(
			"uploading album art: %w",
			apperrors.NewErrBadRequest("album art must be square between 512x512 and 4096x4096"),
		)
	}

	objectKey := uuid.New()
	if err := s.storage.PutObject(
		ctx,
		constants.AlbumArtBucket,
		objectKey.String(),
		reader,
		reader.Size(),
		storage.PutOptions{ContentType: contentType},
	); err != nil {
		return fmt.Errorf("uploading album art: %w", err)
	}

	if err := s.repo.UpdateAlbum(ctx, &UpdateAlbumRequest{
		ID:          albumID,
		AlbumArtKey: &objectKey,
	}); err != nil {
		if errObj := s.storage.DeleteObject(
			context.Background(),
			constants.AlbumArtBucket,
			objectKey.String(),
			storage.DeleteOptions{},
		); errObj != nil {
			return fmt.Errorf(
				"uploading album art: %w, attempting to compensate via saga: %w",
				err,
				errObj,
			)
		}
		return fmt.Errorf("uploading album art: %w", err)
	}

	return nil
}
