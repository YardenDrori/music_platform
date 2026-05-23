package user

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenPair struct {
	accessToken        string
	rawRefreshToken    string
	hashedRefreshToken string
	refreshDur         time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
}

type JWTTokenizerHS256 struct {
	signingKey         []byte
	accessTokenDurSex  time.Duration
	refreshTokenDurSec time.Duration
	hasher             tokenHasher
}

func NewJwtTokenizer(
	signingkey []byte,
	accessTokenDurSex time.Duration,
	refreshTokenDurSec time.Duration,
) *JWTTokenizerHS256 {
	return &JWTTokenizerHS256{
		signingKey:         signingkey,
		accessTokenDurSex:  accessTokenDurSex,
		refreshTokenDurSec: refreshTokenDurSec,
		hasher:             &sha256TokenHasher{},
	}
}

func (t *JWTTokenizerHS256) newAccess(
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

func (t *JWTTokenizerHS256) GenerateTokenPair(
	user *User,
) (*tokenPair, error) {
	accessToken, err := t.newAccess(user)
	if err != nil {
		return nil, fmt.Errorf("generating token pair: %w", err)
	}

	randString := rand.Text() + rand.Text()
	rawRefreshToken := randString[:32]
	hashedRefreshToken := t.hasher.hashToken(rawRefreshToken)

	return &tokenPair{
		accessToken:        *accessToken,
		rawRefreshToken:    rawRefreshToken,
		hashedRefreshToken: hashedRefreshToken,
		refreshDur:         t.refreshTokenDurSec,
	}, nil
}

func (t *JWTTokenizerHS256) ValidateAccessToken(
	ctx context.Context,
	token string,
) (*Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(token, claims,
		func(tok *jwt.Token) (any, error) {
			if tok.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid signing method")
			}
			return t.signingKey, nil
		})

	switch {
	case err == nil:
		return claims, nil
	case errors.Is(err, jwt.ErrTokenExpired):
		return nil, ErrExpiredToken
	default:
		return nil, fmt.Errorf("validating token: %w", err)
	}
}
