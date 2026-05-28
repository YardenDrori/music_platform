package user

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockRepo struct{}

func (r *mockRepo) Create(ctx context.Context, u *User) error {
	return nil
}

func (r *mockRepo) Update(ctx context.Context, u *User) error {
	return nil
}

func (r *mockRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *mockRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	return &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		Username:     "testuser",
		FirstName:    "Test",
		LastName:     "User",
		PasswordHash: strings.Repeat("a", 32),
		CreatedAt:    time.Now(),
		LastUpdated:  time.Now(),
	}, nil
}

func (r *mockRepo) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {

	return &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		Username:     "testuser",
		FirstName:    "Test",
		LastName:     "User",
		PasswordHash: strings.Repeat("a", 32),
		CreatedAt:    time.Now(),
		LastUpdated:  time.Now(),
	}, nil
}

func (r *mockRepo) FindByUsername(ctx context.Context, username string) (*User, error) {
	return &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		Username:     "testuser",
		FirstName:    "Test",
		LastName:     "User",
		PasswordHash: strings.Repeat("a", 32),
		CreatedAt:    time.Now(),
		LastUpdated:  time.Now(),
	}, nil
}

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name    string
		req     *NewUserRequest
		wantErr bool
	}{{
		name: "valid request",
		req: &NewUserRequest{
			ID:        uuid.New(),
			Email:     "valid@email.com",
			Username:  "testuser",
			FirstName: "Test",
			LastName:  "User",
			Password:  "superdupersecretpassword",
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(&mockRepo{}, NewArgon2idPasswordHasher())
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
