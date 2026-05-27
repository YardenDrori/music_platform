package user

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

type mockRepo struct{}




func TestNewAccount(t *testing.T) {
	tests := []struct {
		name    string
		req     *NewUserRequest
		wantErr bool
	}{
		name: "valid request",
		req: &NewUserRequest{
			ID: uuid.New(),
			Email: "valid@email.com",
			Username: "testuser",
			FirstName: "Test",
			LastName: "User",
			Password: "superdupersecretpassword",
		},
		wantErr: false,
	},

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})
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
