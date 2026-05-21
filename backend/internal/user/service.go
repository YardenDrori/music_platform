package user

import (
	"context"
)

type AuthService interface {
	//errors:
	//[ErrBadRequest]
	//[ErrUnauthorized]
	Register(ctx context.Context, req *registerRequest) (*authServiceResponse, error)

	//errors:
	//[ErrBadRequest]
	//[ErrUnauthorized]
	//[fmt.Errorf]
	Login(ctx context.Context, req *loginRequest) (*authServiceResponse, error)
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
