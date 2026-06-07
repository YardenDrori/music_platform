package artists

import (
	"context"
	"fmt"
	"net/url"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func validateArtist(artist *Artist) error {
	if artist.Name == "" {
		return &apperrors.ErrBadRequest{Message: "artist name cannot be empty"}
	}
	if artist.BirthPlace != nil && utf8.RuneCountInString(*artist.BirthPlace) != 2 {
		return &apperrors.ErrBadRequest{Message: "birth place must be 2 characters"}
	}
	if artist.BirthDate != nil && artist.BirthDate.After(time.Now()) {
		return &apperrors.ErrBadRequest{Message: "birth date cannot be in the future"}
	}

	if artist.LinkToYouTube != nil {
		u, err := url.Parse(*artist.LinkToYouTube)
		if err != nil || u.Host == "" {
			return &apperrors.ErrBadRequest{Message: "link to youtube is not a valid url"}
		}
	}
	if artist.LinkToSpotify != nil {
		u, err := url.Parse(*artist.LinkToSpotify)
		if err != nil || u.Host == "" {
			return &apperrors.ErrBadRequest{Message: "link to spotify is not a valid url"}
		}
	}
	if artist.LinkToAppleMusic != nil {
		u, err := url.Parse(*artist.LinkToAppleMusic)
		if err != nil || u.Host == "" {
			return &apperrors.ErrBadRequest{Message: "link to apple music is not a valid url"}
		}
	}

	return nil
}

func validateUpdateReq(req *UpdateArtistReq) error {
	if req.Name != nil && *req.Name == "" {
		return &apperrors.ErrBadRequest{Message: "artist name cannot be empty"}
	}
	if req.BirthPlace != nil && utf8.RuneCountInString(*req.BirthPlace) != 2 {
		return &apperrors.ErrBadRequest{Message: "birth place must be 2 characters"}
	}
	if req.BirthDate != nil && req.BirthDate.After(time.Now()) {
		return &apperrors.ErrBadRequest{Message: "birth date cannot be in the future"}
	}
	if req.LinkToYouTube != nil {
		u, err := url.Parse(*req.LinkToYouTube)
		if err != nil || u.Host == "" {
			return &apperrors.ErrBadRequest{Message: "link to youtube is not a valid url"}
		}
	}
	if req.LinkToSpotify != nil {
		u, err := url.Parse(*req.LinkToSpotify)
		if err != nil || u.Host == "" {
			return &apperrors.ErrBadRequest{Message: "link to spotify is not a valid url"}
		}
	}
	if req.LinkToAppleMusic != nil {
		u, err := url.Parse(*req.LinkToAppleMusic)
		if err != nil || u.Host == "" {
			return &apperrors.ErrBadRequest{Message: "link to apple music is not a valid url"}
		}
	}
	return nil
}

func (s *service) NewArtist(ctx context.Context, req NewArtistReq) error {
	artist := req.ToArtist()
	if err := validateArtist(&artist); err != nil {
		return fmt.Errorf("creating artist: %w", err)
	}

	err := s.repo.NewArtist(ctx, artist)
	if err != nil {
		return fmt.Errorf("creating artist: %w", err)
	}
	return nil
}

func (s *service) GetArtistsByName(ctx context.Context, name string) ([]*Artist, error) {
	artists, err := s.repo.GetArtistsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("getting artists by name: %w", err)
	}
	return artists, nil
}

func (s *service) GetArtistByID(ctx context.Context, id uuid.UUID) (*Artist, error) {
	artist, err := s.repo.GetArtistByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting artist by id: %w", err)
	}
	return artist, nil
}

func (s *service) UpdateArtistDetails(ctx context.Context, req UpdateArtistReq) error {
	if err := validateUpdateReq(&req); err != nil {
		return fmt.Errorf("updating artist: %w", err)
	}
	if err := s.repo.UpdateArtist(ctx, req); err != nil {
		return fmt.Errorf("updating artist: %w", err)
	}
	return nil
}

func (s *service) SoftDeleteArtist(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	err := s.repo.UpdateArtist(ctx, UpdateArtistReq{ID: id, DeletedAt: &now})
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
