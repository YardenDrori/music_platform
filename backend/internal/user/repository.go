package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	Alter(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByEmail(ctx context.Context, e string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, u string) (*User, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, email, username, 
		 password_hash, created_at, is_active)
	   VALUES($1, $2, $3, $4, $5, $6)`,
		u.ID, u.Email, u.Username, u.PasswordHash, u.CreatedAt, true,
	)
	return err
}
func (r *postgresRepository) Alter(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET email = $1, username = $2, password_hash = $3, is_active = $4 WHERE id = $5`,
		u.Email, u.Username, u.PasswordHash, u.Active, u.ID,
	)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *postgresRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *postgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, password_hash, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *postgresRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, username, password_hash, created_at FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
