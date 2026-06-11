package artists

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (r *postgresRepository) NewArtist(
	ctx context.Context,
	artist Artist,
	uploaderID uuid.UUID,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("creating new artist: %w", apperrors.NewErrInternal().WithCause(err))
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO artists (id, name,  description, is_band, artist_image_key,
		artist_banner_key, link_to_youtube, link_to_spotify, link_to_apple_music,
		origin_date, origin_place, added_at, updated_at, deleted_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		artist.ID,
		artist.Name,
		artist.Description,
		artist.IsBand,
		artist.ArtistImageKey,
		artist.ArtistBannerKey,
		artist.LinkToYouTube,
		artist.LinkToSpotify,
		artist.LinkToAppleMusic,
		artist.OriginDate,
		artist.OriginPlace,
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
		return fmt.Errorf("creating new artist: %w", apperrors.NewErrInternal().WithCause(err))
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO artist_contributors (artist_id, user_id, contributed_at) VALUES ($1, $2, $3)`,
		artist.ID,
		uploaderID,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("creating new artist: %w", apperrors.NewErrInternal().WithCause(err))
	}

	for _, alias := range artist.Aliases {
		_, err = tx.Exec(ctx,
			`INSERT INTO artist_aliases (artist_id, alias) VALUES($1, $2)`,
			artist.ID,
			alias,
		)
		if err != nil {
			return fmt.Errorf("creating new artist: %w", apperrors.NewErrInternal().WithCause(err))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("creating new artist: %w", apperrors.NewErrInternal().WithCause(err))
	}
	return nil
}

func (r *postgresRepository) GetArtistsByNameOrAlias(
	ctx context.Context,
	name string,
) ([]ArtistEntity, error) {
	rows, err := r.db.Query(ctx, `
	SELECT a.id, a.name, a.description, a.is_band, a.artist_image_key, a.artist_banner_key,
	a.link_to_youtube, a.link_to_spotify, a.link_to_apple_music, a.origin_date, a.origin_place,
	a.added_at, a.updated_at, 
	ARRAY(SELECT alias FROM artist_aliases WHERE artist_id = a.id ORDER BY alias),
	ARRAY(SELECT user_id FROM artist_contributors ac WHERE ac.artist_id = a.id ORDER BY ac.user_id),
	ARRAY(SELECT username FROM users u JOIN artist_contributors ac ON ac.user_id = u.id WHERE a.id = ac.artist_id ORDER BY ac.user_id),
	ARRAY(SELECT profile_pic_key FROM users u JOIN artist_contributors ac ON ac.user_id = u.id WHERE a.id = ac.artist_id ORDER BY ac.user_id),
	ARRAY(SELECT contributed_at FROM artist_contributors ac WHERE ac.artist_id = a.id ORDER BY ac.user_id)
	FROM artists a WHERE a.deleted_at IS NULL
	AND (a.name ILIKE '%' || $1 || '%' OR EXISTS (
		SELECT 1 FROM artist_aliases aa WHERE aa.artist_id = a.id AND aa.alias ILIKE '%' || $1 || '%'
	))`,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"fetching artists via name or alias: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}
	defer rows.Close()

	artists := []ArtistEntity{}
	for rows.Next() {
		artist := ArtistEntity{}
		aliases := []string{}

		userIDs := []uuid.UUID{}
		userNames := []string{}
		userProfilePics := []*uuid.UUID{}
		contributionDates := []time.Time{}
		err := rows.Scan(
			&artist.ID,
			&artist.Name,
			&artist.Description,
			&artist.IsBand,
			&artist.ArtistImageKey,
			&artist.ArtistBannerKey,
			&artist.LinkToYouTube,
			&artist.LinkToSpotify,
			&artist.LinkToAppleMusic,
			&artist.OriginDate,
			&artist.OriginPlace,
			&artist.AddedAt,
			&artist.UpdatedAt,
			&aliases,
			&userIDs,
			&userNames,
			&userProfilePics,
			&contributionDates,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"fetching artists via name or alias: %w",
				apperrors.NewErrInternal().WithCause(err),
			)
		}

		contributions := []ContributionEntity{}
		//all arrays are the same length
		for i := range userIDs {
			contribution := ContributionEntity{
				ContributorID:         userIDs[i],
				ContributorName:       userNames[i],
				ContributorProfileKey: userProfilePics[i],
				ContributionDate:      contributionDates[i],
			}
			contributions = append(contributions, contribution)
		}

		artist.Aliases = aliases
		artist.ContributionsEntity = contributions
		artists = append(artists, artist)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf(
			"fetching artists via name or alias: %w",
			apperrors.NewErrInternal().WithCause(rows.Err()),
		)
	}
	if len(artists) == 0 {
		return nil, fmt.Errorf(
			"fetching artists via name or alias: %w",
			apperrors.NewErrNotFound("artist not found"),
		)
	}
	return artists, nil
}

func (r *postgresRepository) GetArtistByID(
	ctx context.Context,
	id uuid.UUID,
) (*ArtistEntity, error) {
	artist := ArtistEntity{}
	aliases := []string{}

	userIDs := []uuid.UUID{}
	userNames := []string{}
	userProfilePics := []*uuid.UUID{}
	contributionDates := []time.Time{}
	err := r.db.QueryRow(ctx, `
	SELECT a.id, a.name, a.description, a.is_band, a.artist_image_key, a.artist_banner_key,
	a.link_to_youtube, a.link_to_spotify, a.link_to_apple_music, a.origin_date, a.origin_place,
	a.added_at, a.updated_at, 
	ARRAY(SELECT alias FROM artist_aliases WHERE artist_id = a.id ORDER BY alias),
	ARRAY(SELECT user_id FROM artist_contributors ac WHERE ac.artist_id = a.id ORDER BY ac.user_id),
	ARRAY(SELECT username FROM users u JOIN artist_contributors ac ON ac.user_id = u.id WHERE a.id = ac.artist_id ORDER BY ac.user_id),
	ARRAY(SELECT profile_pic_key FROM users u JOIN artist_contributors ac ON ac.user_id = u.id WHERE a.id = ac.artist_id ORDER BY ac.user_id),
	ARRAY(SELECT contributed_at FROM artist_contributors ac WHERE ac.artist_id = a.id ORDER BY ac.user_id)
	FROM artists a WHERE a.deleted_at IS NULL AND a.id = $1`,
		id,
	).Scan(
		&artist.ID,
		&artist.Name,
		&artist.Description,
		&artist.IsBand,
		&artist.ArtistImageKey,
		&artist.ArtistBannerKey,
		&artist.LinkToYouTube,
		&artist.LinkToSpotify,
		&artist.LinkToAppleMusic,
		&artist.OriginDate,
		&artist.OriginPlace,
		&artist.AddedAt,
		&artist.UpdatedAt,
		&aliases,
		&userIDs,
		&userNames,
		&userProfilePics,
		&contributionDates,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf(
			"fetching artists via id: %w",
			apperrors.NewErrNotFound("artist not found"),
		)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"fetching artists via id: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}

	contributions := []ContributionEntity{}
	//all arrays are the same length
	for i := range userIDs {
		contribution := ContributionEntity{
			ContributorID:         userIDs[i],
			ContributorName:       userNames[i],
			ContributorProfileKey: userProfilePics[i],
			ContributionDate:      contributionDates[i],
		}
		contributions = append(contributions, contribution)
	}
	artist.Aliases = aliases
	artist.ContributionsEntity = contributions

	return &artist, nil
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
	if req.IsBand != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_band = $%d", i))
		args = append(args, *req.IsBand)
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
	if req.OriginDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("origin_date = $%d", i))
		args = append(args, *req.OriginDate)
		i++
	}
	if req.OriginPlace != nil {
		setClauses = append(setClauses, fmt.Sprintf("origin_place = $%d", i))
		args = append(args, *req.OriginPlace)
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
		return fmt.Errorf("updating artist: %w", apperrors.NewErrInternal().WithCause(err))
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("artist not found")
	}
	return nil
}

func (r *postgresRepository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	data, err := r.db.Exec(ctx, `DELETE FROM artists WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting artist: %w", apperrors.NewErrInternal().WithCause(err))
	}
	if data.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("artist not found")
	}
	return nil
}

func (r *postgresRepository) AddContributor(
	ctx context.Context,
	artistID uuid.UUID,
	userID uuid.UUID,
) error {
	_, err := r.db.Exec(ctx, `INSERT INTO artist_contributors (artist_id, user_id) VALUES ($1, $2)`,
		artistID, userID,
	)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return apperrors.NewErrConflict("user is already a contributor")
			}
		}
		return fmt.Errorf("adding contributor: %w", apperrors.NewErrInternal().WithCause(err))
	}
	return nil
}

func (r *postgresRepository) RemoveContributor(
	ctx context.Context,
	artistID uuid.UUID,
	userID uuid.UUID,
) error {
	tag, err := r.db.Exec(
		ctx,
		`DELETE FROM artist_contributors WHERE artist_id = $1 AND user_id = $2`,
		artistID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("removing contributor: %w", apperrors.NewErrInternal().WithCause(err))
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("contributor not found")
	}
	return nil
}

func (r *postgresRepository) AddAlias(
	ctx context.Context,
	artistID uuid.UUID,
	alias string,
) error {
	_, err := r.db.Exec(ctx, `INSERT INTO artist_aliases (artist_id, alias) VALUES($1, $2)`,
		artistID, alias,
	)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return apperrors.NewErrConflict("artist already has this alias")
			}
		}
		return fmt.Errorf("adding alias: %w", apperrors.NewErrInternal().WithCause(err))
	}
	return nil
}

func (r *postgresRepository) RemoveAlias(
	ctx context.Context,
	artistID uuid.UUID,
	alias string,
) error {
	tag, err := r.db.Exec(
		ctx,
		`DELETE FROM artist_aliases WHERE artist_id = $1 AND alias = $2`,
		artistID,
		alias,
	)
	if err != nil {
		return fmt.Errorf("removing alias: %w", apperrors.NewErrInternal().WithCause(err))
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("alias not found for this artist")
	}
	return nil
}
