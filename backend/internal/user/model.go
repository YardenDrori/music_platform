package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	Email         string
	Username      string
	FirstName     string
	LastName      string
	ProfilePicKey *uuid.UUID
	PasswordHash  string
	CreatedAt     time.Time
	LastUpdated   time.Time
	// Active        bool
}

type NewUserRequest struct {
	ID        uuid.UUID
	Email     string
	Username  string
	FirstName string
	LastName  string
	Password  string
}
