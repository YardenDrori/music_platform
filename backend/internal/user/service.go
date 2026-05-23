package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByEmail(ctx context.Context, e string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, u string) (*User, error)
	NewToken(
		ctx context.Context,
		id uuid.UUID,
		tokenHash string,
		iat time.Time,
		exp time.Time,
	) error
	DeleteToken(ctx context.Context, id uuid.UUID, tokenHash string) error
	FindToken(ctx context.Context, tokenHash string) (*uuid.UUID, error)
	CleanExpiredTokens(ctx context.Context) error
}

type TokenValidator interface {
	ValidateAccessToken(ctx context.Context, token string) (*Claims, error)
}

type tokenizer interface {
	TokenValidator
	GenerateTokenPair(user *User) (*tokenPair, error)
}

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
	tokenizer tokenizer
}

func NewService(repo Repository, tok tokenizer) AuthService {
	return &userService{
		repo:      repo,
		hasher:    &argon2idPasswordHasher{},
		tokenizer: tok,
	}
}
