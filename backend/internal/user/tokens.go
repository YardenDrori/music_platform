package user

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

type Tokenizer interface {
	ValidateAccessToken(ctx context.Context, token string) (Claims, error)
}

type JWTTokenizerHS256 struct {
	signingKey         []byte
	accessTokenDurSex  time.Duration
	refreshTokenDurSec time.Duration
}

func NewJwtTokenizer(
	signingkey []byte,
	accessTokenDurSex time.Duration,
	refreshTokenDurSec time.Duration,
) Tokenizer {
	return &JWTTokenizerHS256{
		signingKey:         signingkey,
		accessTokenDurSex:  accessTokenDurSex,
		refreshTokenDurSec: refreshTokenDurSec,
	}
}

func (t *JWTTokenizerHS256) newAccess(
	ctx context.Context,
	user *User,
) (*string, error) {
	expTime := time.Now().UTC().Add(t.accessTokenDurSex)

	claims := Claims{
		jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	finalizedToken, err := token.SignedString(t.signingKey)
	if err != nil {
		return nil, fmt.Errorf("generating new access token: %w", err)
	}
	return &finalizedToken, nil
}

func (t *JWTTokenizerHS256) generateTokenPair(
	ctx context.Context,
	user *User,
) (string, string, error) {
	panic("")
}

func (t *JWTTokenizerHS256) ValidateAccessToken(ctx context.Context, token string) (Claims, error) {
	panic("")
}
