package albums

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

type repository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) handleQueryError(
	currentlyDoing string,
	conflictOn string,
	err error,
) error {
	if pgerr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgerr.Code {
		case "23505":
			return fmt.Errorf(
				"%v: %w",
				currentlyDoing,
				apperrors.NewErrConflict(fmt.Sprintf("%v already exists", conflictOn)))
		case "23503":
			return fmt.Errorf(
				"%v: %w",
				currentlyDoing,
				apperrors.NewErrBadRequest(fmt.Sprintf("%v does not exist", conflictOn)))
		}
	}
	return fmt.Errorf(
		"%v: %w",
		currentlyDoing,
		apperrors.NewErrInternal().WithCause(err),
	)
}

func (r *repository) NewAlbum(ctx context.Context, album *Album) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("creating new album: %w", apperrors.NewErrInternal().WithCause(err))
	}
	//nolint
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
	INSERT INTO albums (id, name, description, main_artist_id, album_art_key,
	has_all_tracks, added_at, updated_at, deleted_at, premiered_at)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		album.ID,
		album.Name,
		album.Description,
		album.MainArtistID,
		album.AlbumArtKey,
		album.HasAllTracks,
		album.AddedAt,
		album.UpdatedAt,
		album.DeletedAt,
		album.PremieredAt,
	)
	if err != nil {
		return r.handleQueryError("creating new album", "album", err)
	}

	for _, artist := range album.Artists {
		_, err = tx.Exec(ctx, `
			INSERT INTO album_artists (album_id, artist_id) VALUES($1, $2)`,
			album.ID,
			artist,
		)
		if err != nil {
			return r.handleQueryError("creating new album", "artist for this album", err)
		}
	}

	for _, contributor := range album.Contributors {
		_, err = tx.Exec(ctx, `
			INSERT INTO album_contributors (album_id, user_id) VALUES($1, $2) ON CONFLICT DO NOTHING`,
			album.ID,
			contributor.ContributorID,
		)
		if err != nil {
			return r.handleQueryError(
				"creating new album",
				"contributor for this album",
				err,
			)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("creating new album: %w", apperrors.NewErrInternal().WithCause(err))
	}

	return nil
}

func (r *repository) GetAlbumsByName(ctx context.Context, name string) ([]Album, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.name, a.description, a.main_artist_id, a.album_art_key,
		a.has_all_tracks, a.added_at, a.updated_at, a.deleted_at, a.premiered_at,
		ARRAY(SELECT artist_id FROM album_artists aa WHERE aa.album_id = a.id),
		ARRAY(SELECT user_id FROM album_contributors ac WHERE ac.album_id = a.id ORDER BY ac.user_id),
		ARRAY(SELECT username FROM users u JOIN album_contributors ac ON u.id = ac.user_id WHERE ac.album_id = a.id ORDER BY ac.user_id),
		ARRAY(SELECT profile_pic_key FROM users u JOIN album_contributors ac ON u.id = ac.user_id WHERE ac.album_id = a.id ORDER BY ac.user_id),
		ARRAY(SELECT contributed_at FROM album_contributors ac WHERE ac.album_id = a.id ORDER BY ac.user_id)
		FROM albums a
		WHERE deleted_at IS NULL AND name ILIKE '%' || $1 || '%'
		ORDER BY name ASC`,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"getting albums via name: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}
	defer rows.Close()

	albums := []Album{}
	for rows.Next() {
		album := Album{}
		artistIDs := []uuid.UUID{}
		userIDs := []uuid.UUID{}
		profilePicKeys := []*string{}
		userNames := []string{}
		contributionDates := []time.Time{}

		err := rows.Scan(
			&album.ID,
			&album.Name,
			&album.Description,
			&album.MainArtistID,
			&album.AlbumArtKey,
			&album.HasAllTracks,
			&album.AddedAt,
			&album.UpdatedAt,
			&album.DeletedAt,
			&album.PremieredAt,
			&artistIDs,
			&userIDs,
			&userNames,
			&profilePicKeys,
			&contributionDates,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"fetching albums via name: %w",
				apperrors.NewErrInternal().WithCause(err),
			)
		}

		album.Artists = append(album.Artists, artistIDs...)

		//all contribution related arrays are the same length
		for i := range userIDs {
			newContributor := Contributor{
				ContributorID:         userIDs[i],
				ContributorName:       userNames[i],
				ContributorProfileUrl: profilePicKeys[i],
				ContributionDate:      contributionDates[i],
			}
			album.Contributors = append(album.Contributors, newContributor)
		}

		albums = append(albums, album)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf(
			"fetching albums via name: %w",
			apperrors.NewErrInternal().WithCause(rows.Err()),
		)
	}
	if len(albums) == 0 {
		return nil, fmt.Errorf(
			"fetching albums via name: %w",
			apperrors.NewErrNotFound("album not found"),
		)
	}

	return albums, nil
}

func (r *repository) GetAlbumByID(ctx context.Context, id uuid.UUID) (*Album, error) {
	album := Album{}
	artistIDs := []uuid.UUID{}
	userIDs := []uuid.UUID{}
	profilePicKeys := []*string{}
	userNames := []string{}
	contributionDates := []time.Time{}

	err := r.db.QueryRow(ctx, `
		SELECT a.id, a.name, a.description, a.main_artist_id, a.album_art_key,
		a.has_all_tracks, a.added_at, a.updated_at, a.deleted_at, a.premiered_at,
		ARRAY(SELECT artist_id FROM album_artists aa WHERE aa.album_id = a.id),
		ARRAY(SELECT user_id FROM album_contributors ac WHERE ac.album_id = a.id ORDER BY ac.user_id),
		ARRAY(SELECT username FROM users u JOIN album_contributors ac ON u.id = ac.user_id WHERE ac.album_id = a.id ORDER BY ac.user_id),
		ARRAY(SELECT profile_pic_key FROM users u JOIN album_contributors ac ON u.id = ac.user_id WHERE ac.album_id = a.id ORDER BY ac.user_id),
		ARRAY(SELECT contributed_at FROM album_contributors ac WHERE ac.album_id = a.id ORDER BY ac.user_id)
		FROM albums a
		WHERE deleted_at IS NULL AND a.id = $1`,
		id,
	).Scan(
		&album.ID,
		&album.Name,
		&album.Description,
		&album.MainArtistID,
		&album.AlbumArtKey,
		&album.HasAllTracks,
		&album.AddedAt,
		&album.UpdatedAt,
		&album.DeletedAt,
		&album.PremieredAt,
		&artistIDs,
		&userIDs,
		&userNames,
		&profilePicKeys,
		&contributionDates,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(
				"getting album via id: %w",
				apperrors.NewErrNotFound("album not found"),
			)
		}
		return nil, fmt.Errorf(
			"getting album via id: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}

	album.Artists = append(album.Artists, artistIDs...)

	//all contribution related arrays are the same length
	for i := range userIDs {
		newContributor := Contributor{
			ContributorID:         userIDs[i],
			ContributorName:       userNames[i],
			ContributorProfileUrl: profilePicKeys[i],
			ContributionDate:      contributionDates[i],
		}
		album.Contributors = append(album.Contributors, newContributor)
	}

	return &album, nil
}

func (r *repository) UpdateAlbum(ctx context.Context, req *UpdateAlbumRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"updating album information: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}
	//nolint
	defer tx.Rollback(ctx)

	setClauses := []string{"updated_at = NOW()"}
	args := []any{}

	if req.Name != nil {
		args = append(args, *req.Name)
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(args)))
	}
	if req.Description != nil {
		args = append(args, *req.Description)
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", len(args)))
	}
	if req.MainArtistID != nil {
		args = append(args, *req.MainArtistID)
		setClauses = append(setClauses, fmt.Sprintf("main_artist_id = $%d", len(args)))
	}
	if req.AlbumArtKey != nil {
		args = append(args, *req.AlbumArtKey)
		setClauses = append(setClauses, fmt.Sprintf("album_art_key = $%d", len(args)))
	}
	if req.HasAllTracks != nil {
		args = append(args, *req.HasAllTracks)
		setClauses = append(setClauses, fmt.Sprintf("has_all_tracks = $%d", len(args)))
	}
	if req.DeletedAt != nil {
		args = append(args, *req.DeletedAt)
		setClauses = append(setClauses, fmt.Sprintf("deleted_at = $%d", len(args)))
	}
	if req.PremieredAt != nil {
		args = append(args, *req.PremieredAt)
		setClauses = append(setClauses, fmt.Sprintf("premiered_at = $%d", len(args)))
	}
	args = append(args, req.ID)

	tag, err := tx.Exec(ctx, fmt.Sprintf("UPDATE albums SET %v WHERE id = $%d",
		strings.Join(setClauses, ", "),
		len(args),
	),
		args...,
	)
	if err != nil {
		return r.handleQueryError("updating album information", "main artist", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf(
			"updating album information: %w",
			apperrors.NewErrNotFound("album not found"),
		)
	}

	for _, artist := range req.ArtistsToAdd {
		_, err := tx.Exec(ctx, `
		INSERT INTO album_artists (album_id, artist_id) VALUES($1, $2)`,
			req.ID,
			artist,
		)
		if err != nil {
			return r.handleQueryError(
				"updating album information",
				"artist for this album",
				err,
			)
		}
	}
	for _, artist := range req.ArtistsToRemove {
		tag, err := tx.Exec(ctx, `
		DELETE FROM album_artists WHERE album_id = $1 AND artist_id = $2`,
			req.ID,
			artist,
		)
		if err != nil {
			return fmt.Errorf(
				"updating album information: %w",
				apperrors.NewErrInternal().WithCause(err),
			)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf(
				"updating album information: %w",
				apperrors.NewErrNotFound("artist isn't featured in this album"),
			)
		}
	}
	for _, contr := range req.ContributorIDsToAdd {
		_, err := tx.Exec(ctx, `
		INSERT INTO album_contributors (album_id, user_id) VALUES($1, $2) ON CONFLICT DO NOTHING`,
			req.ID,
			contr,
		)
		if err != nil {
			return r.handleQueryError(
				"updating album information",
				"contributor for this album",
				err,
			)
		}
	}
	for _, contr := range req.ContributorIDsToRemove {
		tag, err := tx.Exec(ctx, `
		DELETE FROM album_contributors WHERE album_id = $1 AND user_id = $2`,
			req.ID,
			contr,
		)
		if err != nil {
			return fmt.Errorf(
				"updating album information: %w",
				apperrors.NewErrInternal().WithCause(err),
			)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf(
				"updating album information: %w",
				apperrors.NewErrNotFound("user isn't a contributor of this album"),
			)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"updating album information: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}
	return nil
}

func (r *repository) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM albums WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("deleting album: %w", apperrors.NewErrInternal().WithCause(err))
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("deleting album: %w", apperrors.NewErrNotFound("album not found"))
	}
	return nil
}
