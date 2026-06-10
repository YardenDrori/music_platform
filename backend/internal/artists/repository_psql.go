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

func (r *postgresRepository) NewArtist(ctx context.Context, artist Artist, uploaderID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("creating new artist: %w", apperrors.NewErrInternal("").WithCause(err))
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
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
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return apperrors.NewErrConflict("conflict")
			}
		}
		return fmt.Errorf("creating new artist: %w", apperrors.NewErrInternal("").WithCause(err))
	}

	_, err = tx.Exec(ctx, `INSERT INTO artist_contributors (artist_id, user_id) VALUES ($1, $2)`,
		artist.ID, uploaderID,
	)
	if err != nil {
		return fmt.Errorf("creating new artist: %w", apperrors.NewErrInternal("").WithCause(err))
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("creating new artist: %w", apperrors.NewErrInternal("").WithCause(err))
	}
	return nil
}

func (r *postgresRepository) fetchContributors(ctx context.Context, artistID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `SELECT user_id FROM artist_contributors WHERE artist_id = $1`, artistID)
	if err != nil {
		return nil, fmt.Errorf("fetching contributors: %w", apperrors.NewErrInternal("").WithCause(err))
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("fetching contributors: %w", apperrors.NewErrInternal("").WithCause(err))
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetching contributors: %w", apperrors.NewErrInternal("").WithCause(err))
	}
	return ids, nil
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
		return nil, fmt.Errorf("finding artists by name: %w", apperrors.NewErrInternal("").WithCause(err))
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
			return nil, fmt.Errorf("finding artists by name: %w", apperrors.NewErrInternal("").WithCause(err))
		}
		newArtist.ContributorIDs, err = r.fetchContributors(ctx, newArtist.ID)
		if err != nil {
			return nil, fmt.Errorf("finding artists by name: %w", err)
		}
		artists = append(artists, newArtist)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("finding artists by name: %w", apperrors.NewErrInternal("").WithCause(err))
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
		birth_place, added_at, updated_at, deleted_at FROM artists WHERE id = $1`,
		id,
	).Scan(&artist.ID, &artist.Name, &artist.Description, &artist.ArtistImageKey,
		&artist.ArtistBannerKey, &artist.LinkToYouTube, &artist.LinkToSpotify,
		&artist.LinkToAppleMusic, &artist.BirthDate, &artist.BirthPlace,
		&artist.AddedAt, &artist.UpdatedAt, &artist.DeletedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewErrNotFound("artist not found")
		}
		return nil, fmt.Errorf("finding artist by id: %w", apperrors.NewErrInternal("").WithCause(err))
	}

	artist.ContributorIDs, err = r.fetchContributors(ctx, artist.ID)
	if err != nil {
		return nil, fmt.Errorf("finding artist by id: %w", err)
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
	if req.ArtistImageKey != nil {
		setClauses = append(setClauses, fmt.Sprintf("artist_image_key = $%d", i))
		args = append(args, *req.ArtistImageKey)
		i++
	}
	if req.ArtistBannerKey != nil {
		setClauses = append(setClauses, fmt.Sprintf("artist_banner_key = $%d", i))
		args = append(args, *req.ArtistBannerKey)
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
		return fmt.Errorf("updating artist: %w", apperrors.NewErrInternal("").WithCause(err))
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("artist not found")
	}
	return nil
}

func (r *postgresRepository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	data, err := r.db.Exec(ctx, `DELETE FROM artists WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting artist: %w", apperrors.NewErrInternal("").WithCause(err))
	}
	if data.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("artist not found")
	}
	return nil
}

func (r *postgresRepository) AddContributor(ctx context.Context, artistID uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `INSERT INTO artist_contributors (artist_id, user_id) VALUES ($1, $2)`,
		artistID, userID,
	)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return apperrors.NewErrConflict("user is already a contributor")
			}
		}
		return fmt.Errorf("adding contributor: %w", apperrors.NewErrInternal("").WithCause(err))
	}
	return nil
}

func (r *postgresRepository) RemoveContributor(ctx context.Context, artistID uuid.UUID, userID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM artist_contributors WHERE artist_id = $1 AND user_id = $2`,
		artistID, userID,
	)
	if err != nil {
		return fmt.Errorf("removing contributor: %w", apperrors.NewErrInternal("").WithCause(err))
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("contributor not found")
	}
	return nil
}
