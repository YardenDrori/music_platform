package user

import (
	"context"
	"fmt"
	"go/token"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
	ID       uuid.UUID
	Username string
	Email    string
}

type Tokenizer interface {
	ValidateAccessToken(ctx context.Context, token string) (Claims, error)
}

type JWTTokenizer struct {
	signingKey []byte
}

func NewJwtTokenizer(signingkey []byte) Tokenizer {
	return &JWTTokenizer{signingKey: signingkey}
}

func (t *JWTTokenizer) newAccess(
	ctx context.Context,
	user *User,
) (string, error) {
	panic("")
}

func (t *JWTTokenizer) generateTokenPair(
	ctx context.Context,
	user *User,
) (string, string, error) {
	panic("")
}

func (t *JWTTokenizer) ValidateAccessToken(ctx context.Context, token string) (Claims, error) {
	panic("")
}
