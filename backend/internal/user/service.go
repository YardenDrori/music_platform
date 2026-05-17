package user

import (
	"context"
)

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	Login(ctx context.Context, req *LoginRequest) (*User, error)
}

type userService struct {
	repo   Repository
	hasher passwordHasher
}

func NewService(repo Repository) AuthService {
	return &userService{repo: repo, hasher: &argon2idPasswordHasher{}}
}
