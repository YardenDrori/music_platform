package user

import (
	"context"

	"github.com/google/uuid"
)

type repository interface {
	Create(ctx context.Context, u *User) error
	// errors:
	// ErrConflict,
	// errorf
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

	// errors:
	// [ErrUnathenticated]
	// [ErrForbidden]
	// [ErrNotFound]
	// [errorf]
	FindByUUID(ctx context.Context, id uuid.UUID) (*User, error)

	// errors:
	// ErrBadRequest,
	// ErrConflict
	// ErrForbidden
	// ErrUnathenticated
	// errorf
	UpdateAccount(ctx context.Context, user *User) error

	DeleteAccount(ctx context.Context, id uuid.UUID) error

	// errors:
	// ErrUnathenticated,
	// ErrForbidden,
	// ErrConflict,
	// errorf
	DeactivateAccount(ctx context.Context, id uuid.UUID) error
}
