package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	FirstName    string
	LastName     string
	PasswordHash string
	CreatedAt    time.Time
	Active       bool
}
