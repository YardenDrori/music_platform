package user

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
)

type mockRepo struct {
	createFn         func(ctx context.Context, u *User) error
	findByEmailFn    func(ctx context.Context, email string) (*User, error)
	findByUsernameFn func(ctx context.Context, username string) (*User, error)
	findByIDFn       func(ctx context.Context, id uuid.UUID) (*User, error)
}

func (r *mockRepo) Create(ctx context.Context, u *User) error {
	if r.createFn != nil {
		return r.createFn(ctx, u)
	}
	return nil
}

func (r *mockRepo) Update(ctx context.Context, u *User) error {
	return nil
}

func (r *mockRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *mockRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	if r.findByEmailFn != nil {
		return r.findByEmailFn(ctx, email)
	}
	return nil, apperrors.ErrNotFound
}

func (r *mockRepo) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if r.findByIDFn != nil {
		return r.findByIDFn(ctx, id)
	}
	return nil, apperrors.ErrNotFound
}

func (r *mockRepo) FindByUsername(ctx context.Context, username string) (*User, error) {
	if r.findByUsernameFn != nil {
		return r.findByUsernameFn(ctx, username)
	}
	return nil, apperrors.ErrNotFound
}

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name    string
		req     *NewUserRequest
		repo    *mockRepo
		wantErr bool
	}{
		{
			name: "happy",
			req: &NewUserRequest{
				ID:        uuid.New(),
				Email:     "valid@email.com",
				Username:  "testuser",
				FirstName: "Test",
				LastName:  "User",
				Password:  "superdupersecretpassword",
			},
			repo:    &mockRepo{},
			wantErr: false,
		},
		{
			name: "db err on create",
			req: &NewUserRequest{
				ID:        uuid.New(),
				Email:     "valid@email.com",
				Username:  "testuser",
				FirstName: "Test",
				LastName:  "User",
				Password:  "superdupersecretpassword",
			},
			repo: &mockRepo{
				createFn: func(ctx context.Context, u *User) error {
					return fmt.Errorf("db is literally on fire HELP!")
				},
			},
			wantErr: true,
		},
		{
			name: "password too short",
			req: &NewUserRequest{
				ID:        uuid.New(),
				Email:     "valid@email.com",
				Username:  "testuser",
				FirstName: "Test",
				LastName:  "User",
				Password:  "shorty",
			},
			repo:    &mockRepo{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.repo, NewArgon2idPasswordHasher())
			_, err := service.NewAccount(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("got err = %v, wantErr = %v, request = %v", err, tt.wantErr, tt.req)
			}
		})
	}
}

func TestValidateAccountBusinessRules(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &User{
				Email:        "test@example.com",
				Username:     "testuser",
				FirstName:    "Test",
				LastName:     "User",
				PasswordHash: strings.Repeat("a", 32),
			},
			wantErr: false,
		},
		{
			name: "missing email",
			user: &User{
				Username:     "testuser",
				FirstName:    "Test",
				LastName:     "User",
				PasswordHash: strings.Repeat("a", 32),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAccountBusinessRules(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("got err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
