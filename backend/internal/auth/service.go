package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/user"
)

type repository interface {
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
	// errors:
	// [ErrBadToken]
	// [ErrBadRequest]
	// [Errorf]
	ValidateAccessToken(ctx context.Context, token string) (*Claims, error)
}

type tokenizer interface {
	TokenValidator
	GenerateTokenPair(userId uuid.UUID) (*tokenPair, error)
}

type tokenHasher interface {
	hashToken(token string) string
}

type Service interface {
	TokenValidator

	//errors:
	//[ErrBadRequest]
	//[ErrUnauthorized]
	Register(ctx context.Context, req *registerRequest) (*authServiceResponse, *user.User, error)

	//errors:
	//[ErrBadRequest]
	//[ErrUnauthorized]
	//[fmt.Errorf]
	Login(ctx context.Context, req *loginRequest) (*authServiceResponse, *user.User, error)

	//errors:
	//[ErrBadRequest]
	//[fmt.Errorf]
	RequestAccessToken(ctx context.Context, oldRawToken string) (*authServiceResponse, error)
}
