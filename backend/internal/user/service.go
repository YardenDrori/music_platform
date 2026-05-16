package user

import (
	"context"
)

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest)
}

type userService struct {
	repo Repository
}

func NewService(repo Repository) *userService {
	return &userService{repo: repo}
}
