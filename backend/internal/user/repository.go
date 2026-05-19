package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	//account actions
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByEmail(ctx context.Context, e string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, u string) (*User, error)

	//token actions
	NewToken(
		ctx context.Context,
		id uuid.UUID,
		tokenHash string,
		iat time.Time,
		exp time.Time,
	) error
	DeleteToken(ctx context.Context, id uuid.UUID, tokenHash string) error
	FindToken(ctx context.Context, id uuid.UUID, tokenHash string) (bool, error)
	CleanExpiredTokens(ctx context.Context) error
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

// ==========ACCOUNTS==========
func (r *postgresRepository) Create(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, email, username,
			first_name, last_name,
			password_hash, created_at, is_active)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		u.ID, u.Email, u.Username, u.FirstName, u.LastName,
		u.PasswordHash, u.CreatedAt, u.Active,
	)

	if err == nil {
		return nil
	}
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == "23505" {
			return ErrConflict
		}
	}

	return fmt.Errorf("creating new user in postgres db: %w", err)
}

func (r *postgresRepository) Update(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET email = $1, username = $2, first_name = $3, last_name = $4,
		password_hash = $5, is_active = $6 WHERE id = $7`,
		u.Email, u.Username, u.FirstName, u.LastName, u.PasswordHash, u.Active, u.ID,
	)

	if err == nil {
		return nil
	}
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == "23505" {
			return ErrConflict
		}
	}

	return fmt.Errorf("updating user: %w", err)
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
		password_hash, created_at, is_active FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName,
		&u.PasswordHash, &u.CreatedAt, &u.Active)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("finding user by email: %w", err)
	}

	return u, nil
}

func (r *postgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, first_name, last_name,
		password_hash, created_at, is_active FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName,
		&u.PasswordHash, &u.CreatedAt, &u.Active)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("finding user by uuid: %w", err)
	}

	return u, nil
}

func (r *postgresRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, first_name, last_name,
		password_hash, created_at, is_active FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName,
		&u.PasswordHash, &u.CreatedAt, &u.Active)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("finding user by username: %w", err)
	}

	return u, nil
}

// ==========TOKENS==========
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
			return ErrConflict
		}
	}
	return fmt.Errorf("creating token user in postgres db: %w", err)
}

func (r *postgresRepository) FindToken(
	ctx context.Context,
	id uuid.UUID,
	tokenHash string,
) (bool, error) {
	var dummy int
	err := r.db.QueryRow(ctx,
		`SELECT 1 FROM tokens WHERE token_hash = $1 AND user_id = $2 AND expires_at > NOW()
	`, tokenHash, id).Scan(&dummy)

	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("verifying token against db: %w", err)
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
