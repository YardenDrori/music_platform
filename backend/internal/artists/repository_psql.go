package artists

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
		birth_date, birth_place, added_at, updated_at, deleted_at, uploader_id) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
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
		artist.UploaderID,
	)

	if err == nil {
		return nil
	}
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == "23505" {
			return apperrors.NewErrConflict("conflict")
		}
	}

	return fmt.Errorf("creating new artist in postgres db: %w", err)
}

func (r *postgresRepository) GetArtistsByName(ctx context.Context, name string) ([]*Artist, error) {
	var artists []*Artist
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description, artist_image_key, artist_banner_key,
		link_to_youtube, link_to_spotify, link_to_apple_music, birth_date,
		birth_place, added_at, updated_at, deleted_at, uploader_id FROM artists WHERE name = $1`,
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
			&newArtist.UploaderID,
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
		return nil, apperrors.NewErrNotFound("artist not found")
	}
	return artists, nil
}

func (r *postgresRepository) GetArtistByID(ctx context.Context, id uuid.UUID) (*Artist, error) {
	artist := &Artist{}
	err := r.db.QueryRow(ctx, `
		SELECT id, name, description, artist_image_key, artist_banner_key,
		link_to_youtube, link_to_spotify, link_to_apple_music, birth_date,
		birth_place, added_at, updated_at, deleted_at, uploader_id FROM artists WHERE id = $1`,
		id,
	).Scan(&artist.ID, &artist.Name, &artist.Description, &artist.ArtistImageKey,
		&artist.ArtistBannerKey, &artist.LinkToYouTube, &artist.LinkToSpotify,
		&artist.LinkToAppleMusic, &artist.BirthDate, &artist.BirthPlace,
		&artist.AddedAt, &artist.UpdatedAt, &artist.DeletedAt, &artist.UploaderID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewErrNotFound("artist not found")
		}
		return nil, fmt.Errorf("finding artist by uuid: %w", err)
	}
	return artist, nil
}

func (r *postgresRepository) UpdateArtist(ctx context.Context, req *UpdateArtistReq) error {
	setClauses := []string{"updated_at = NOW()"}
	args := []any{}
	i := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", i))
		args = append(args, *req.Name)
		i++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", i))
		args = append(args, *req.Description)
		i++
	}
	if req.LinkToYouTube != nil {
		setClauses = append(setClauses, fmt.Sprintf("link_to_youtube = $%d", i))
		args = append(args, *req.LinkToYouTube)
		i++
	}
	if req.LinkToSpotify != nil {
		setClauses = append(setClauses, fmt.Sprintf("link_to_spotify = $%d", i))
		args = append(args, *req.LinkToSpotify)
		i++
	}
	if req.LinkToAppleMusic != nil {
		setClauses = append(setClauses, fmt.Sprintf("link_to_apple_music = $%d", i))
		args = append(args, *req.LinkToAppleMusic)
		i++
	}
	if req.BirthDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("birth_date = $%d", i))
		args = append(args, *req.BirthDate)
		i++
	}
	if req.BirthPlace != nil {
		setClauses = append(setClauses, fmt.Sprintf("birth_place = $%d", i))
		args = append(args, *req.BirthPlace)
		i++
	}
	if req.DeletedAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("deleted_at = $%d", i))
		args = append(args, *req.DeletedAt)
		i++
	}

	args = append(args, req.ID)
	query := fmt.Sprintf("UPDATE artists SET %s WHERE id = $%d", strings.Join(setClauses, ", "), i)

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("updating artist: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("artist not found")
	}
	return nil
}

func (r *postgresRepository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	data, err := r.db.Exec(ctx, `DELETE FROM artists WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting artist: %w", err)
	}
	if data.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("artist not found")
	}
	return nil
}
