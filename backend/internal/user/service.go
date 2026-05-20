package user

import (
	"context"
)

type AuthService interface {
	Register(ctx context.Context, req *registerRequest) (*User, error)
	Login(ctx context.Context, req *loginRequest) (*authResponse, error)
}

type userService struct {
	repo      Repository
	hasher    passwordHasher
	tokenizer Tokenizer
}

func NewService(repo Repository, tokenizer Tokenizer) AuthService {
	return &userService{
		repo:      repo,
		hasher:    &argon2idPasswordHasher{},
		tokenizer: tokenizer,
	}
}
