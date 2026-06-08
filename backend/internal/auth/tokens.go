package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
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
	signingKey      []byte
	accessTokenDur  time.Duration
	refreshTokenDur time.Duration
	hasher          tokenHasher
}

func NewJwtTokenizer(
	signingkey []byte,
	accessTokenDur time.Duration,
	refreshTokenDur time.Duration,
) *JWTTokenizerHS256 {
	return &JWTTokenizerHS256{
		signingKey:      signingkey,
		accessTokenDur:  accessTokenDur,
		refreshTokenDur: refreshTokenDur,
		hasher:          &sha256TokenHasher{},
	}
}

func (t *JWTTokenizerHS256) newAccess(
	id uuid.UUID,
) (*string, error) {
	expTime := time.Now().UTC().Add(t.accessTokenDur)

	claims := Claims{
		jwt.RegisteredClaims{
			Subject:   id.String(),
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
	userID uuid.UUID,
) (*tokenPair, error) {
	accessToken, err := t.newAccess(userID)
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
		refreshDur:         t.refreshTokenDur,
	}, nil
}

func (t *JWTTokenizerHS256) ValidateAccessToken(
	ctx context.Context,
	token string,
) (*Claims, error) {
	if token == "" {
		return nil, apperrors.NewErrBadRequest("token not provided")
	}

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims,
		func(tok *jwt.Token) (any, error) {
			if tok.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid signing method: %s", tok.Method)
			}
			return t.signingKey, nil
		})
	if err != nil {
		return nil, fmt.Errorf(
			"validating access token: %w",
			apperrors.NewErrBadToken("token invalid or expired").WithCause(err),
		)
	}
	return claims, nil
}
