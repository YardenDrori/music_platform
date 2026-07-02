package songs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	"github.com/YardenDrori/music-platform/internal/summaries"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func handleQueryErrors(err error, operationDesc string, operationOn string) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		//conflict
		case "23505":
			return fmt.Errorf(
				"%v: %w",
				operationDesc,
				apperrors.NewErrConflict(fmt.Sprintf("%v already exists", operationOn)))
		//constraint violation
		case "23503":
			return fmt.Errorf(
				"%v: %w",
				operationDesc,
				apperrors.NewErrBadRequest(fmt.Sprintf("%v does not exist", operationOn)))
		}
	}
	return fmt.Errorf(
		"%v: %w",
		operationDesc,
		apperrors.NewErrInternal().WithCause(err),
	)
}

func (r *repository) NewSong(ctx context.Context, song NewSong) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("registering new song: %w", apperrors.NewErrInternal().WithCause(err))
	}
	//nolint
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		ctx,
		`INSERT INTO songs (id, title, main_artist_id, album_id, track_number, premiered_at, runtime_ms,
		song_type, upload_method, is_public, audio_key, cover_art_key, added_at, updated_at, deleted_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		song.ID,
		song.Title,
		song.MainArtistID,
		song.AlbumID,
		song.TrackNumber,
		song.PremieredAt,
		song.Runtime.Milliseconds(),
		song.SongType,
		song.UploadMethod,
		song.IsPublic,
		song.AudioKey,
		song.CoverArtKey,
		song.AddedAt,
		song.UpdatedAt,
		song.DeletedAt,
	)
	if err != nil {
		return handleQueryErrors(err, "registering new song", "song")
	}

	//main artist is a "featured artist"
	if song.MainArtistID != nil {
		if _, err = tx.Exec(ctx, `INSERT INTO song_artists (song_id, artist_id) VALUES($1, $2)`,
			song.ID, song.MainArtistID); err != nil {
			return handleQueryErrors(err, "registering new song", "song artist")
		}
	}
	for i := range song.ArtistIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO song_artists (song_id, artist_id) VALUES($1, $2)`,
			song.ID, song.ArtistIDs[i]); err != nil {
			return handleQueryErrors(err, "registering new song", "song artist")
		}
	}

	for i := range song.Contributions {
		if _, err := tx.Exec(ctx, `INSERT INTO song_contributors (song_id, user_id, contributed_at)
		VALUES($1, $2, $3)`,
			song.ID,
			song.Contributions[i].UserSummary.ID,
			song.Contributions[i].ContributionDate,
		); err != nil {
			return handleQueryErrors(err, "registering new song", "song contribution")
		}
	}

	for i := range song.WhitelistedUserIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO song_whitelisted_users (song_id, user_id)
		VALUES($1, $2)`,
			song.ID,
			song.WhitelistedUserIDs[i],
		); err != nil {
			return handleQueryErrors(err, "registering new song", "whitelisted user")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("registering new song: %w", apperrors.NewErrInternal().WithCause(err))
	}

	if _, err := r.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY active_songs`); err != nil {
		slog.Error("registering new song: failed to refresh view", "error", err)
	}

	return nil
}

// querier lets the get-by-id helper run against either the pool (r.db) or an
// open transaction (tx), so Update/Delete can lock the row with FOR UPDATE.
type querier interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

func getSongByID(
	ctx context.Context,
	q querier,
	id uuid.UUID,
	forUpdate bool,
) (*SongRow, error) {
	query := `
		SELECT id, title, main_artist_id, album_id, track_number,
		premiered_at, runtime_ms, song_type, upload_method, is_public,
		audio_key, cover_art_key, added_at, updated_at, deleted_at`
	if forUpdate {
		// writes hit the raw table (materialized view isn't lockable) and skip
		// already soft-deleted rows.
		query += ` FROM songs WHERE deleted_at IS NULL AND id = $1 FOR UPDATE`
	} else {
		query += ` FROM active_songs WHERE id = $1`
	}

	song := &SongRow{}
	var runtimeMs int64

	err := q.QueryRow(ctx, query, id).Scan(
		&song.ID, &song.Title, &song.MainArtistID, &song.AlbumID, &song.TrackNumber,
		&song.PremieredAt, &runtimeMs, &song.SongType, &song.UploadMethod, &song.IsPublic,
		&song.AudioKey, &song.CoverArtKey, &song.AddedAt, &song.UpdatedAt, &song.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewErrNotFound("song not found")
		}
		return nil, fmt.Errorf("finding song by id: %w", apperrors.NewErrInternal().WithCause(err))
	}

	// runtime_ms is stored as int milliseconds go's duration is nanoseconds.
	song.Runtime = time.Duration(runtimeMs) * time.Millisecond

	return song, nil
}

func (r *repository) GetSongByID(ctx context.Context, id uuid.UUID) (*SongRow, error) {
	return getSongByID(ctx, r.db, id, false)
}

func (r *repository) GetSongSummariesByName(
	ctx context.Context,
	songName string,
	limit int,
) ([]summaries.SongSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, runtime_ms, is_public, cover_art_key
		FROM active_songs
		WHERE title ILIKE '%' || $1 || '%'
		ORDER BY title
		LIMIT $2`,
		songName,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"fetching song summaries via name: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}
	defer rows.Close()

	songs := []summaries.SongSummary{}
	for rows.Next() {
		summary := summaries.SongSummary{}
		var runtimeMs int64
		if err := rows.Scan(
			&summary.ID, &summary.Title, &runtimeMs, &summary.IsPublic, &summary.CoverArtKey,
		); err != nil {
			return nil, fmt.Errorf(
				"fetching song summaries via name: %w",
				apperrors.NewErrInternal().WithCause(err),
			)
		}
		summary.Runtime = time.Duration(runtimeMs) * time.Millisecond
		songs = append(songs, summary)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf(
			"fetching song summaries via name: %w",
			apperrors.NewErrInternal().WithCause(rows.Err()),
		)
	}

	return songs, nil
}

func (r *repository) UpdateSong(ctx context.Context, req *UpdateSongReq) (*SongRow, error) {
	setClauses := []string{"updated_at = NOW()"}
	args := []any{}
	i := 1

	if req.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", i))
		args = append(args, *req.Title)
		i++
	}
	if req.MainArtistID != nil {
		setClauses = append(setClauses, fmt.Sprintf("main_artist_id = $%d", i))
		args = append(args, *req.MainArtistID)
		i++
	}
	if req.AlbumID != nil {
		setClauses = append(setClauses, fmt.Sprintf("album_id = $%d", i))
		args = append(args, *req.AlbumID)
		i++
	}
	if req.TrackNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("track_number = $%d", i))
		args = append(args, *req.TrackNumber)
		i++
	}
	if req.PremieredAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("premiered_at = $%d", i))
		args = append(args, *req.PremieredAt)
		i++
	}
	if req.SongType != nil {
		setClauses = append(setClauses, fmt.Sprintf("song_type = $%d", i))
		args = append(args, *req.SongType)
		i++
	}
	if req.IsPublic != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_public = $%d", i))
		args = append(args, *req.IsPublic)
		i++
	}
	if req.CoverArtKey != nil {
		setClauses = append(setClauses, fmt.Sprintf("cover_art_key = $%d", i))
		args = append(args, *req.CoverArtKey)
		i++
	}
	if req.DeletedAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("deleted_at = $%d", i))
		args = append(args, *req.DeletedAt)
		i++
	}
	args = append(args, req.ID)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating song: %w", apperrors.NewErrInternal().WithCause(err))
	}
	//nolint
	defer tx.Rollback(context.Background())

	oldSong, err := getSongByID(ctx, tx, req.ID, true)
	if err != nil {
		return nil, fmt.Errorf("updating song: fetching old song: %w", err)
	}

	tag, err := tx.Exec(ctx, fmt.Sprintf("UPDATE songs SET %v WHERE id = $%d",
		strings.Join(setClauses, ", "),
		i),
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("updating song: %w", apperrors.NewErrInternal().WithCause(err))
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("updating song: %w", apperrors.NewErrNotFound("song not found"))
	}

	for _, artistID := range req.ArtistsToAdd {
		if _, err := tx.Exec(ctx,
			`INSERT INTO song_artists (song_id, artist_id) VALUES($1, $2) ON CONFLICT DO NOTHING`,
			req.ID, artistID,
		); err != nil {
			return nil, handleQueryErrors(err, "updating song", "song artist")
		}
	}
	for _, artistID := range req.ArtistsToRemove {
		tag, err := tx.Exec(ctx,
			`DELETE FROM song_artists WHERE song_id = $1 AND artist_id = $2`,
			req.ID, artistID,
		)
		if err != nil {
			return nil, fmt.Errorf("updating song: %w", apperrors.NewErrInternal().WithCause(err))
		}
		if tag.RowsAffected() == 0 {
			return nil, fmt.Errorf("updating song: %w", apperrors.NewErrNotFound("song artist not found"))
		}
	}

	for _, contr := range req.ContributorsToAdd {
		if _, err := tx.Exec(ctx,
			`INSERT INTO song_contributors (song_id, user_id, contributed_at) VALUES($1, $2, $3)
			ON CONFLICT DO NOTHING`,
			req.ID, contr.UserSummary.ID, contr.ContributionDate,
		); err != nil {
			return nil, handleQueryErrors(err, "updating song", "song contribution")
		}
	}
	for _, userID := range req.ContributorsToRemove {
		tag, err := tx.Exec(ctx,
			`DELETE FROM song_contributors WHERE song_id = $1 AND user_id = $2`,
			req.ID, userID,
		)
		if err != nil {
			return nil, fmt.Errorf("updating song: %w", apperrors.NewErrInternal().WithCause(err))
		}
		if tag.RowsAffected() == 0 {
			return nil, fmt.Errorf("updating song: %w", apperrors.NewErrNotFound("song contribution not found"))
		}
	}

	for _, userID := range req.WhitelistToAdd {
		if _, err := tx.Exec(ctx,
			`INSERT INTO song_whitelisted_users (song_id, user_id) VALUES($1, $2) ON CONFLICT DO NOTHING`,
			req.ID, userID,
		); err != nil {
			return nil, handleQueryErrors(err, "updating song", "whitelisted user")
		}
	}
	for _, userID := range req.WhitelistToRemove {
		tag, err := tx.Exec(ctx,
			`DELETE FROM song_whitelisted_users WHERE song_id = $1 AND user_id = $2`,
			req.ID, userID,
		)
		if err != nil {
			return nil, fmt.Errorf("updating song: %w", apperrors.NewErrInternal().WithCause(err))
		}
		if tag.RowsAffected() == 0 {
			return nil, fmt.Errorf("updating song: %w", apperrors.NewErrNotFound("whitelisted user not found"))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("updating song: %w", apperrors.NewErrInternal().WithCause(err))
	}

	if _, err := r.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY active_songs`); err != nil {
		slog.Error("updating song: failed to refresh view", "error", err)
	}

	return oldSong, nil
}

func (r *repository) DeleteSong(ctx context.Context, id uuid.UUID) (*SongRow, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("deleting song: %w", apperrors.NewErrInternal().WithCause(err))
	}
	//nolint
	defer tx.Rollback(context.Background())

	oldSong, err := getSongByID(ctx, tx, id, true)
	if err != nil {
		return nil, fmt.Errorf("deleting song: %w", err)
	}

	tag, err := tx.Exec(ctx, `DELETE FROM songs WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("deleting song: %w", apperrors.NewErrInternal().WithCause(err))
	}
	if tag.RowsAffected() == 0 {
		return nil, apperrors.NewErrNotFound("song not found")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("deleting song: %w", apperrors.NewErrInternal().WithCause(err))
	}

	if _, err := r.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY active_songs`); err != nil {
		slog.Error("deleting song: failed to refresh view", "error", err)
	}

	return oldSong, nil
}
