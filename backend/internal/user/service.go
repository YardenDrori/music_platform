package user

import (
	"context"

	"github.com/google/uuid"
)

type repository interface {
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByEmail(ctx context.Context, e string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, u string) (*User, error)
}

type Service interface {
	NewAccount(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}
