package artists

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/url"
	"time"
	"unicode/utf8"

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

func validateArtist(artist *Artist) error {
	if artist.Name == "" {
		return apperrors.NewErrBadRequest("artist name cannot be empty")
	}
	if artist.OriginPlace != nil && utf8.RuneCountInString(*artist.OriginPlace) != 2 {
		return apperrors.NewErrBadRequest("origin place must be 2 characters")
	}
	if artist.OriginDate != nil && artist.OriginDate.After(time.Now()) {
		return apperrors.NewErrBadRequest("origin date cannot be in the future")
	}

	if artist.LinkToYouTube != nil {
		u, err := url.Parse(*artist.LinkToYouTube)
		if err != nil || u.Host == "" {
			return apperrors.NewErrBadRequest("link to youtube is not a valid url")
		}
	}
	if artist.LinkToSpotify != nil {
		u, err := url.Parse(*artist.LinkToSpotify)
		if err != nil || u.Host == "" {
			return apperrors.NewErrBadRequest("link to spotify is not a valid url")
		}
	}
	if artist.LinkToAppleMusic != nil {
		u, err := url.Parse(*artist.LinkToAppleMusic)
		if err != nil || u.Host == "" {
			return apperrors.NewErrBadRequest("link to apple music is not a valid url")
		}
	}

	return nil
}

func validateUpdateReq(req *UpdateArtistReq) error {
	if req.Name != nil && *req.Name == "" {
		return apperrors.NewErrBadRequest("artist name cannot be empty")
	}
	if req.OriginPlace != nil && utf8.RuneCountInString(*req.OriginPlace) != 2 {
		return apperrors.NewErrBadRequest("origin place must be 2 characters")
	}
	if req.OriginDate != nil && req.OriginDate.After(time.Now()) {
		return apperrors.NewErrBadRequest("origin date cannot be in the future")
	}
	if req.LinkToYouTube != nil {
		u, err := url.Parse(*req.LinkToYouTube)
		if err != nil || u.Host == "" {
			return apperrors.NewErrBadRequest("link to youtube is not a valid url")
		}
	}
	if req.LinkToSpotify != nil {
		u, err := url.Parse(*req.LinkToSpotify)
		if err != nil || u.Host == "" {
			return apperrors.NewErrBadRequest("link to spotify is not a valid url")
		}
	}
	if req.LinkToAppleMusic != nil {
		u, err := url.Parse(*req.LinkToAppleMusic)
		if err != nil || u.Host == "" {
			return apperrors.NewErrBadRequest("link to apple music is not a valid url")
		}
	}
	return nil
}

func (s *service) NewArtist(ctx context.Context, req NewArtistReq) error {
	requesterID, err := identity.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("registering new artist: %w", err)
	}

	artist := req.ToArtist()
	if err := validateArtist(&artist); err != nil {
		return fmt.Errorf("creating artist: %w", err)
	}

	err = s.repo.NewArtist(ctx, artist, requesterID)
	if err != nil {
		return fmt.Errorf("creating artist: %w", err)
	}
	return nil
}

func (s *service) processArtistEntity(artist *Artist) {
	if artist.ArtistImageUrl != nil {
		artist.ArtistImageUrl = s.storage.BuildPublicGetUrl(
			constants.ProfilePicBucket,
			*artist.ArtistImageUrl,
		)
	}
	if artist.ArtistBannerUrl != nil {
		artist.ArtistBannerUrl = s.storage.BuildPublicGetUrl(
			constants.BannerBucket,
			*artist.ArtistBannerUrl,
		)
	}
	for j := range artist.Contributions {
		contribution := &artist.Contributions[j]
		if contribution.ContributorProfileUrl != nil {
			contribution.ContributorProfileUrl = s.storage.BuildPublicGetUrl(
				constants.ProfilePicBucket,
				*contribution.ContributorProfileUrl,
			)
		}
	}
}

func (s *service) GetArtistsByNameOrAlias(ctx context.Context, name string) ([]Artist, error) {
	artists, err := s.repo.GetArtistsByNameOrAlias(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("getting artists by name: %w", err)
	}

	for i := range artists {
		s.processArtistEntity(&artists[i])
	}

	return artists, nil
}

func (s *service) GetArtistByID(ctx context.Context, id uuid.UUID) (*Artist, error) {
	artist, err := s.repo.GetArtistByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting artist by id: %w", err)
	}

	s.processArtistEntity(artist)

	return artist, nil
}

func (s *service) UpdateArtistDetails(ctx context.Context, req *UpdateArtistReq) error {
	requesterID, err := identity.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("updating artist details: %w", err)
	}

	if err := validateUpdateReq(req); err != nil {
		return fmt.Errorf("updating artist: %w", err)
	}

	req.ContributorsToAdd = append(req.ContributorsToAdd, requesterID)

	if err := s.repo.UpdateArtist(ctx, req); err != nil {
		return fmt.Errorf("updating artist: %w", err)
	}
	return nil
}

func (s *service) SoftDeleteArtist(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	err := s.repo.UpdateArtist(ctx, &UpdateArtistReq{ID: id, DeletedAt: &now})
	if err != nil {
		return fmt.Errorf("soft deleting artist: %w", err)
	}
	return nil
}

func (s *service) HardDeleteArtist(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteArtist(ctx, id); err != nil {
		return fmt.Errorf("hard deleting artist: %w", err)
	}
	return nil
}

var allowedImageFormats = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
}

func (s *service) UploadArtistProfilePicture(
	ctx context.Context,
	file []byte,
	artistID uuid.UUID,
) error {
	contentType := http.DetectContentType(file)

	if _, ok := allowedImageFormats[contentType]; !ok {
		return fmt.Errorf(
			"uploading artist profile picture: %w",
			apperrors.NewErrBadRequest("invalid content type"),
		)
	}

	reader := bytes.NewReader(file)
	config, _, err := image.DecodeConfig(reader)
	if err != nil {
		return fmt.Errorf(
			"uploading artist profile picture: %w",
			apperrors.NewErrBadRequest("invalid content type"),
		)
	}
	reader.Reset(file)

	if config.Height != 1024 || config.Width != 1024 {
		return fmt.Errorf(
			"uploading artist profile picture: %w",
			apperrors.NewErrBadRequest("image resolution must be 1024x1024"),
		)
	}

	objectKey := uuid.New()
	storageOpts := storage.PutOptions{
		ContentType:    contentType,
		SendContentMD5: false,
	}

	if err := s.storage.PutObject(
		ctx,
		constants.ProfilePicBucket,
		objectKey.String(),
		reader,
		reader.Size(),
		storageOpts,
	); err != nil {
		return fmt.Errorf("uploading artist profile picture: %w", err)
	}
	if err := s.repo.UpdateArtist(ctx, &UpdateArtistReq{
		ID:             artistID,
		ArtistImageKey: &objectKey,
	}); err != nil {
		if errObj := s.storage.DeleteObject(
			context.Background(),
			constants.ProfilePicBucket,
			objectKey.String(),
			storage.DeleteOptions{},
		); errObj != nil {
			return fmt.Errorf(
				"uploading artist profile picture: %w, attempting to remove object via sage: %w",
				err,
				errObj,
			)
		}
		return fmt.Errorf("uploading artist profile picture: %w", err)
	}

	return nil
}

func (s *service) UploadArtistBannerPicture(
	ctx context.Context,
	file []byte,
	artistID uuid.UUID,
) error {
	contentType := http.DetectContentType(file)

	if _, ok := allowedImageFormats[contentType]; !ok {
		return fmt.Errorf(
			"uploading artist banner: %w",
			apperrors.NewErrBadRequest("invalid content type"),
		)
	}

	reader := bytes.NewReader(file)
	config, _, err := image.DecodeConfig(reader)
	if err != nil {
		return fmt.Errorf(
			"uploading artist banner: %w",
			apperrors.NewErrBadRequest("invalid content type"),
		)
	}
	reader.Reset(file)

	if config.Width < 1920 || config.Width > 3840 || config.Width != config.Height*3 {
		return fmt.Errorf(
			"uploading artist banner: %w",
			apperrors.NewErrBadRequest(
				"banner must be 3:1 aspect ratio with width between 1920 and 3840",
			),
		)
	}

	objectKey := uuid.New()
	storageOpts := storage.PutOptions{
		ContentType:    contentType,
		SendContentMD5: false,
	}

	if err := s.storage.PutObject(
		ctx,
		constants.BannerBucket,
		objectKey.String(),
		reader,
		reader.Size(),
		storageOpts,
	); err != nil {
		return fmt.Errorf("uploading artist banner: %w", err)
	}

	if err := s.repo.UpdateArtist(ctx, &UpdateArtistReq{
		ID:              artistID,
		ArtistBannerKey: &objectKey,
	}); err != nil {
		if errObj := s.storage.DeleteObject(
			context.Background(),
			constants.BannerBucket,
			objectKey.String(),
			storage.DeleteOptions{},
		); errObj != nil {
			return fmt.Errorf(
				"uploading artist banner: %w, attempting to remove object via sage: %w",
				err,
				errObj,
			)
		}
		return fmt.Errorf("uploading artist banner: %w", err)
	}

	return nil
}
