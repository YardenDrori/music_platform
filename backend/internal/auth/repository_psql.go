package auth

import (
	"context"
	"errors"
	"fmt"
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

func NewPostgresRepository(db *pgxpool.Pool) repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) NewToken(
	ctx context.Context,
	id uuid.UUID,
	tokenHash string,
	iat time.Time,
	exp time.Time,
) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO tokens (id, issued_at, expires_at, user_id, token_hash)
	VALUES($1, $2, $3, $4, $5)
	`, uuid.New(), iat, exp, id, tokenHash,
	)

	if err == nil {
		return nil
	}
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == "23505" {
			return apperrors.NewErrConflict("already exists").WithCause(pgErr)
		}
	}
	return fmt.Errorf("creating token user in postgres db: %w", err)
}

func (r *postgresRepository) FindToken(
	ctx context.Context,
	tokenHash string,
) (*uuid.UUID, error) {
	var owner uuid.UUID
	err := r.db.QueryRow(ctx,
		`SELECT user_id FROM tokens WHERE token_hash = $1 AND expires_at > NOW()
	`, tokenHash).Scan(&owner)

	if err == nil {
		return &owner, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.NewErrBadToken("token invalid or expired").WithCause(err)
	}
	return nil, fmt.Errorf("verifying token against db: %w", err)
}

func (r *postgresRepository) DeleteToken(
	ctx context.Context,
	id uuid.UUID,
	tokenHash string,
) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM tokens WHERE user_id = $1 AND token_hash = $2
`, id, tokenHash)
	if err == nil {
		return nil
	}
	return fmt.Errorf("deleting token: %w", err)
}

func (r *postgresRepository) CleanExpiredTokens(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `
	DELETE FROM tokens WHERE expires_at < NOW()`)
	if err == nil {
		return nil
	}
	return fmt.Errorf("cleaning expired tokens: %w", err)
}
