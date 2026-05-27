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
	userService user.Service
	repo        repository
	tokenHasher tokenHasher
	tokenizer   tokenizer
}

func NewService(repo repository, tok tokenizer, userService user.Service) *service {
	return &service{
		repo:        repo,
		tokenizer:   tok,
		userService: userService,
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

	newUserReq := user.NewUserRequest{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	tokens, err := s.tokenizer.GenerateTokenPair(newUserReq.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("registering: %w", err)
	}

	var repoErr error
	var newUser *user.User
	for range 3 {
		newUser, repoErr = s.userService.NewAccount(ctx, &newUserReq)
		if !errors.Is(repoErr, ErrConflict) {
			break
		}
		slog.Info("congratulation you should take a lottery ticket now!")
	}
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
			newUserReq.ID,
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
	}, newUser, nil
}

// errors:
// ErrBadRequest
// ErrUnauthenticated
// Errorf
func (s *service) Login(
	ctx context.Context,
	req *loginRequest,
) (*authServiceResponse, *user.User, error) {
	var user *user.User
	var err error

	switch {
	case req.Email != nil && *req.Email != "":
		user, err = s.userService.Authenticate(ctx, *req.Email, req.Password)
	case req.UserName != nil && *req.UserName != "":
		user, err = s.userService.Authenticate(ctx, *req.UserName, req.Password)
	default:
		return nil, nil, &ErrBadRequest{Message: "missing credentials"}
	}
	if errors.Is(err, ErrNotFound) {
		return nil, nil, ErrUnauthenticated
	}
	if err != nil {
		return nil, nil, fmt.Errorf("logging in: %w", err)
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
