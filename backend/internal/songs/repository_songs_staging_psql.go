package songs

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/YardenDrori/music-platform/internal/apperrors"
)

func (r *postgresRepository) NewStagingSong(
	ctx context.Context,
	objectKey string,
	uploadID string,
	userID uuid.UUID,
) error {
	_, err := r.db.Exec(ctx, `INSERT INTO songs_staging (id, object_key, upload_id, uploading_user)
	VALUES($1, $2, $3, $4)`,
		uuid.New(),
		objectKey,
		uploadID,
		userID.String(),
	)
	if err != nil {
		return handleQueryErrors(err, "creating new staged song", "staged song")
	}
	return nil
}

func (r *postgresRepository) StagingSongOwnerFromObjectKey(
	ctx context.Context,
	objectKey string,
) (uuid.UUID, error) {
	userID := uuid.UUID{}
	if err := r.db.QueryRow(ctx, `SELECT uploading_user FROM songs_staging WHERE object_key = $1`, objectKey).
		Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf(
				"fetching staging song owner's id from object key: %w",
				apperrors.NewErrNotFound("song not found").WithInternal("staging song not found"),
			)
		}
		return uuid.Nil, fmt.Errorf("fetching staging song owner's id from object key: %w", err)
	}
	return userID, nil
}

func (r *postgresRepository) DeleteStagingSong(
	ctx context.Context, objectKey string,
) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM songs_staging WHERE object_key = $1`, objectKey)

	if err != nil {
		return fmt.Errorf("deleting song: %w", apperrors.NewErrInternal().WithCause(err))
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewErrNotFound("song not found")
	}
	return nil
}

func (r *postgresRepository) EvictStagingSongs(
	ctx context.Context,
	abort func(objectKey string, uploadID string) error,
) error {
	rows, err := r.db.Query(
		ctx,
		`DELETE FROM songs_staging
		 WHERE added_at < now() - INTERVAL '1 hour'
		 RETURNING object_key, upload_id`,
	)
	if err != nil {
		return fmt.Errorf(
			"evicting expired staging songs: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}

	type evicted struct{ objectKey, uploadID string }
	var pending []evicted
	for rows.Next() {
		var e evicted
		if err := rows.Scan(&e.objectKey, &e.uploadID); err != nil {
			rows.Close()
			return fmt.Errorf(
				"evicting expired staging songs: %w",
				apperrors.NewErrInternal().WithCause(err),
			)
		}
		pending = append(pending, e)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			"evicting expired staging songs: %w",
			apperrors.NewErrInternal().WithCause(err),
		)
	}

	var errs []error
	for _, e := range pending {
		if err := abort(e.objectKey, e.uploadID); err != nil {
			errs = append(errs, fmt.Errorf("aborting upload %s: %w", e.uploadID, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("evicting expired staging songs: %w", errors.Join(errs...))
	}
	return nil
}
