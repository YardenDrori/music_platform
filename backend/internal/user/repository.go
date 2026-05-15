package user

import (
	"context"
	"os/user"

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
	panic("")
}
func (r *postgresRepository) Alter(ctx context.Context, u *User) error {
	panic("")
}
func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	panic("")
}
func (r *postgresRepository) FindByEmail(ctx context.Context, e string) (*User, error) {
	panic("")
}
func (r *postgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	panic("")
}
func (r *postgresRepository) FindByUsername(ctx context.Context, e string) (*User, error) {
	panic("")
}
