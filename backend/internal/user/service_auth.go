package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type authServiceResponse struct {
	authResponse
	refreshToken   string
	refreshExpirey time.Time
}

func (s *userService) Register(
	ctx context.Context,
	req *registerRequest,
) (*authServiceResponse, error) {
	if req.Email == "" || req.FirstName == "" || req.LastName == "" || req.Password == "" ||
		req.UserName == "" {
		return nil, &ErrBadRequest{Message: "missing fields"}
	}

	//email verif
	if before, after, ok := strings.Cut(
		req.Email,
		"@",
	); !ok || after == "" || before == "" ||
		strings.ContainsRune(after, '@') || !strings.Contains(after, ".") {
		return nil, &ErrBadRequest{Message: "invalid email"}
	}

	if utf8.RuneCountInString(req.Password) < 8 {
		return nil, &ErrBadRequest{Message: "password too short"}
	}

	//verify email and username are available
	_, errEmail := s.repo.FindByEmail(ctx, req.Email)
	switch {
	case errEmail == nil:
		return nil, &ErrBadRequest{Message: "email unavailable"}
	case errors.Is(errEmail, ErrNotFound):
		break
	default:
		return nil, fmt.Errorf("checking email availability: %w", errEmail)
	}
	_, errUsername := s.repo.FindByUsername(ctx, req.UserName)
	switch {
	case errUsername == nil:
		return nil, &ErrBadRequest{Message: "username unavailable"}
	case errors.Is(errUsername, ErrNotFound):
		break
	default:
		return nil, fmt.Errorf("checking username availability: %w", errUsername)
	}

	passwordHash := s.hasher.hashPassword(req.Password)

	newUser := User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.UserName,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
		Active:       true,
	}

	tokens, err := s.tokenizer.GenerateTokenPair(&newUser)
	if err != nil {
		return nil, fmt.Errorf("registering: %w", err)
	}

	repoErr := s.repo.Create(ctx, &newUser)
	switch {
	case repoErr == nil:
		break
	case errors.Is(repoErr, ErrConflict):
		return nil, &ErrBadRequest{Message: "email or username unavailable"}
	default:
		return nil, fmt.Errorf("attempting to create new user: %w", repoErr)
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
		return nil, fmt.Errorf("registering: %w", err)
	}

	return &authServiceResponse{
		authResponse: authResponse{
			AccessToken: tokens.accessToken,
			User:        &newUser},
		refreshToken:   tokens.rawRefreshToken,
		refreshExpirey: time.Now().UTC().Add(tokens.refreshDur),
	}, nil
}

func (s *userService) Login(ctx context.Context, req *loginRequest) (*authServiceResponse, error) {
	var user *User
	var err error
	switch {
	case req.Email != nil && *req.Email != "":
		user, err = s.repo.FindByEmail(ctx, *req.Email)
	case req.UserName != nil && *req.UserName != "":
		user, err = s.repo.FindByUsername(ctx, *req.UserName)
	default:
		return nil, &ErrBadRequest{Message: "missing credentials"}
	}
	if errors.Is(err, ErrNotFound) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, fmt.Errorf("logging in: %w", err)
	}

	ok, err := s.hasher.verifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("verifying password: %w", err)
	}
	if !ok {
		return nil, ErrUnauthorized
	}

	//authorized - making tokens
	tokens, err := s.tokenizer.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("logging in: %w", err)
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
		return nil, fmt.Errorf("logging in: %w", err)
	}

	return &authServiceResponse{
		authResponse: authResponse{
			User:        user,
			AccessToken: tokens.accessToken},
		refreshToken:   tokens.rawRefreshToken,
		refreshExpirey: time.Now().UTC().Add(tokens.refreshDur),
	}, nil

}
