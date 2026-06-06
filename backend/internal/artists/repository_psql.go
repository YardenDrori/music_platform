package artists

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/YardenDrori/music-platform/internal/apperrors"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) NewArtist(ctx context.Context, artist Artist) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO artists (id, name, description, artist_image_key,
		artist_banner_key, link_to_youtube, link_to_spotify, link_to_apple_music,
		birth_date, birth_place, added_at, updated_at, deleted_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		artist.ID,
		artist.Name,
		artist.Description,
		artist.ArtistImageKey,
		artist.ArtistBannerKey,
		artist.LinkToYouTube,
		artist.LinkToSpotify,
		artist.LinkToAppleMusic,
		artist.BirthDate,
		artist.BirthPlace,
		artist.AddedAt,
		artist.UpdatedAt,
		artist.DeletedAt,
	)

	if err == nil {
		return nil
	}
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == "23505" {
			return apperrors.ErrConflict
		}
	}

	return fmt.Errorf("creating new artist in postgres db: %w", err)
}

func (r *postgresRepository) GetArtistsByName(ctx context.Context, name string) ([]*Artist, error) {
	var artists []*Artist
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description, artist_image_key, artist_banner_key,
		link_to_youtube, link_to_spotify, link_to_apple_music, birth_date,
		birth_place, added_at, updated_at, deleted_at FROM artists WHERE name = $1`,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("finding artists by name: %w", err)
	}

	defer rows.Close()

	foundOne := false
	for rows.Next() {
		foundOne = true
		newArtist := &Artist{}
		err = rows.Scan(
			&newArtist.ID,
			&newArtist.Name,
			&newArtist.Description,
			&newArtist.ArtistImageKey,
			&newArtist.ArtistBannerKey,
			&newArtist.LinkToYouTube,
			&newArtist.LinkToSpotify,
			&newArtist.LinkToAppleMusic,
			&newArtist.BirthDate,
			&newArtist.BirthPlace,
			&newArtist.AddedAt,
			&newArtist.UpdatedAt,
			&newArtist.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("finding artists by name: %w", err)
		}
		artists = append(artists, newArtist)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("fetching all artists named %v: %w", name, err)
	}
	if !foundOne {
		return nil, apperrors.ErrNotFound
	}
	return artists, nil
}

func (r *postgresRepository) GetArtistByID(ctx context.Context, id uuid.UUID) (*Artist, error) {
	artist := &Artist{}
	err := r.db.QueryRow(ctx, `
		SELECT id, name, description, artist_image_key, artist_banner_key,
		link_to_youtube, link_to_spotify, link_to_apple_music, birth_date,
		birth_place, added_at, updated_at, deleted_at FROM artists WHERE id = $1`,
		id,
	).Scan(&artist.ID, &artist.Name, &artist.Description, &artist.ArtistImageKey,
		&artist.ArtistBannerKey, &artist.LinkToYouTube, &artist.LinkToSpotify,
		&artist.LinkToAppleMusic, &artist.BirthDate, &artist.BirthPlace,
		&artist.AddedAt, &artist.UpdatedAt, &artist.DeletedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("finding artist by uuid: %w", err)
	}
	return artist, nil
}

func (r *postgresRepository) UpdateArtist(ctx context.Context, artist Artist) error {
	queryDetails, err := r.db.Exec(ctx, `
		UPDATE artists SET name = $1, description = $2, artist_image_key = $3,
		artist_banner_key = $4, link_to_youtube = $5, link_to_spotify = $6,
		link_to_apple_music = $7, birth_date = $8, birth_place = $9,
		updated_at = $10, deleted_at = $11 WHERE id = $12`,
		artist.Name,
		artist.Description,
		artist.ArtistImageKey,
		artist.ArtistBannerKey,
		artist.LinkToYouTube,
		artist.LinkToSpotify,
		artist.LinkToAppleMusic,
		artist.BirthDate,
		artist.BirthPlace,
		artist.UpdatedAt,
		artist.DeletedAt,
		artist.ID,
	)
	if err != nil {
		return fmt.Errorf("updating artist: %w", err)
	}
	if queryDetails.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *postgresRepository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	data, err := r.db.Exec(ctx, `DELETE FROM artists WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting artist: %w", err)
	}
	if data.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
