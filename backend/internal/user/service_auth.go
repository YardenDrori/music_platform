package user

import (
	"context"
	"strings"
)

func (s *userService) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	if before, after, ok := strings.Cut(
		req.Email,
		"@",
	); !ok || after == "" || before == "" ||
		strings.ContainsRune(after, '@') || !strings.Contains(after, ".") {
		return nil, &ErrBadRequest{Message: "invalid email"}
	}

	// var userName = strings.Clone(req.Email)

	panic("")
}

func (s *userService) Login(ctx context.Context, req *LoginRequest) (*User, error) {
	panic("")
}
