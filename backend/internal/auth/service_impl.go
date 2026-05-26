package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/user"
)

type service struct {
	userService    user.Service
	repo           repository
	passwordHasher passwordHasher
	tokenHasher    tokenHasher
	tokenizer      tokenizer
}

func NewService(repo repository, tok tokenizer, userService user.Service) *service {
	return &service{
		repo:           repo,
		passwordHasher: &argon2idPasswordHasher{},
		tokenizer:      tok,
		userService:    userService,
	}
}

type authServiceResponse struct {
	accessToken    string
	refreshToken   string
	refreshExpirey time.Time
}

func (s *service) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	return s.tokenizer.ValidateAccessToken(ctx, token)
}

func (s *service) Register(
	ctx context.Context,
	req *registerRequest,
) (*authServiceResponse, *user.User, error) {
	if req.Email == "" || req.FirstName == "" || req.LastName == "" || req.Password == "" ||
		req.UserName == "" {
		return nil, nil, &ErrBadRequest{Message: "missing fields"}
	}

	//email verif
	if before, after, ok := strings.Cut(
		req.Email,
		"@",
	); !ok || after == "" || before == "" ||
		strings.ContainsRune(after, '@') || !strings.Contains(after, ".") {
		return nil, nil, &ErrBadRequest{Message: "invalid email"}
	}

	if utf8.RuneCountInString(req.Password) < 8 {
		return nil, nil, &ErrBadRequest{Message: "password too short"}
	}

	passwordHash := s.passwordHasher.hashPassword(req.Password)

	newUser := user.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.UserName,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
		Active:       true,
	}

	tokens, err := s.tokenizer.GenerateTokenPair(newUser.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("registering: %w", err)
	}

	repoErr := s.userService.NewAccount(ctx, &newUser)
	switch {
	case repoErr == nil:
		break
	case errors.Is(repoErr, ErrConflict):
		return nil, nil, &ErrBadRequest{Message: "email or username unavailable"}
	default:
		return nil, nil, fmt.Errorf("attempting to create new user: %w", repoErr)
	}

	for range 3 {
		err = s.repo.NewToken(
			ctx,
			newUser.ID,
			tokens.hashedRefreshToken,
			time.Now().UTC(),
			time.Now().UTC().Add(tokens.refreshDur),
		)
		if !errors.Is(err, ErrConflict) {
			break
		}
		slog.Info("congratulation you should take a lottery ticket now!")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("registering: %w", err)
	}

	return &authServiceResponse{
		accessToken:    tokens.accessToken,
		refreshToken:   tokens.rawRefreshToken,
		refreshExpirey: time.Now().UTC().Add(tokens.refreshDur),
	}, &newUser, nil
}

func (s *service) Login(
	ctx context.Context,
	req *loginRequest,
) (*authServiceResponse, *user.User, error) {
	var user *user.User
	var err error

	switch {
	case req.Email != nil && *req.Email != "":
		user, err = s.userService.FindByEmail(ctx, *req.Email)
	case req.UserName != nil && *req.UserName != "":
		user, err = s.userService.FindByUsername(ctx, *req.UserName)
	default:
		return nil, nil, &ErrBadRequest{Message: "missing credentials"}
	}

	if errors.Is(err, ErrNotFound) {
		return nil, nil, ErrUnauthorized
	}
	if err != nil {
		return nil, nil, fmt.Errorf("logging in: %w", err)
	}

	ok, err := s.passwordHasher.verifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, nil, fmt.Errorf("verifying password: %w", err)
	}
	if !ok {
		return nil, nil, ErrUnauthorized
	}

	//authorized - making tokens
	tokens, err := s.tokenizer.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("logging in: %w", err)
	}

	for range 3 {
		err = s.repo.NewToken(
			ctx,
			user.ID,
			tokens.hashedRefreshToken,
			time.Now().UTC(),
			time.Now().UTC().Add(tokens.refreshDur),
		)
		if !errors.Is(err, ErrConflict) {
			break
		}
		slog.Info("congratulation you should take a lottery ticket now!")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("logging in: %w", err)
	}

	return &authServiceResponse{
		accessToken:    tokens.accessToken,
		refreshToken:   tokens.rawRefreshToken,
		refreshExpirey: time.Now().UTC().Add(tokens.refreshDur),
	}, user, nil

}

func (s *service) RequestAccessToken(
	ctx context.Context,
	oldRawToken string,
) (*authServiceResponse, error) {
	if oldRawToken == "" {
		return nil, &ErrBadRequest{Message: "token not provided"}
	}

	originalRefreshHash := s.tokenHasher.hashToken(oldRawToken)

	owner, err := s.repo.FindToken(ctx, originalRefreshHash)

	switch {
	case err == nil:
		break
	case errors.Is(err, ErrNotFound):
		return nil, ErrBadToken
	default:
		return nil, fmt.Errorf("finding token: %w", err)
	}

	tokens, err := s.tokenizer.GenerateTokenPair(*owner)
	if err != nil {
		return nil, fmt.Errorf("requesting new access token with refresh token cycling: %w", err)
	}

	err = s.repo.DeleteToken(ctx, *owner, originalRefreshHash)
	if err != nil {
		return nil, fmt.Errorf("requesting new access token with refresh token cycling: %w", err)
	}

	err = s.repo.NewToken(
		ctx,
		*owner,
		tokens.hashedRefreshToken,
		time.Now().UTC(),
		time.Now().UTC().Add(tokens.refreshDur),
	)
	if err != nil {
		return nil, fmt.Errorf("requesting new access token with refresh token cycling: %w", err)
	}

	return &authServiceResponse{
		accessToken:    tokens.accessToken,
		refreshToken:   tokens.rawRefreshToken,
		refreshExpirey: time.Now().UTC().Add(tokens.refreshDur),
	}, nil

}
