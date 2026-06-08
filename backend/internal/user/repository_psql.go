package user

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

func NewPostgresRepository(db *pgxpool.Pool) repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, email, username,
			first_name, last_name,
			password_hash, created_at, is_active)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		u.ID, u.Email, u.Username, u.FirstName, u.LastName,
		u.PasswordHash, u.CreatedAt, u.Active,
	)

	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return apperrors.NewErrConflict("already exists").
					WithInternal(pgErr.Detail).
					WithCause(err)
			}
		}
		return fmt.Errorf("creating new user in postgres db: %w", err)
	}
	return nil
}

// errors:
// ErrConflict,
// errorf
func (r *postgresRepository) Update(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET email = $1, username = $2, first_name = $3, last_name = $4,
		password_hash = $5, is_active = $6 WHERE id = $7`,
		u.Email, u.Username, u.FirstName, u.LastName, u.PasswordHash, u.Active, u.ID,
	)

	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return apperrors.NewErrConflict("already exists").
					WithInternal(pgErr.Detail).
					WithCause(err)
			}
		}
		return fmt.Errorf("creating new user in postgres db: %w", err)
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err == nil {
		return nil
	}
	return fmt.Errorf("deleting user: %w", err)
}

func (r *postgresRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, first_name, last_name,
		password_hash, created_at, last_updated, is_active FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName,
		&u.PasswordHash, &u.CreatedAt, &u.LastUpdated, &u.Active)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewErrNotFound("user not found")
		}
		return nil, fmt.Errorf("finding user by email: %w", err)
	}

	return u, nil
}

// errors:
// [ErrNotFound]
// [errorf]
func (r *postgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, first_name, last_name,
		password_hash, created_at, last_updated, is_active FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName,
		&u.PasswordHash, &u.CreatedAt, &u.LastUpdated, &u.Active)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewErrNotFound("user not found")
		}
		return nil, fmt.Errorf("finding user by uuid: %w", err)
	}

	return u, nil
}

func (r *postgresRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, first_name, last_name,
		password_hash, created_at, last_updated, is_active FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName,
		&u.PasswordHash, &u.CreatedAt, &u.LastUpdated, &u.Active)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewErrNotFound("user not found")
		}
		return nil, fmt.Errorf("finding user by username: %w", err)
	}

	return u, nil
}
