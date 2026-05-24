package user

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"
)

type service struct {
	repo repository
}

// func NewService(repo *repository) Service {
// 	return &service{
// 		repo: repo,
// 	}
// }

func (s *service) NewAccount(ctx context.Context, user *User) error {
	//validation logic
	if utf8.RuneCountInString(user.Username) < 3 || utf8.RuneCountInString(user.Email) < 3 ||
		utf8.RuneCountInString(user.FirstName) < 3 ||
		utf8.RuneCountInString(user.LastName) < 3 ||
		utf8.RuneCountInString(user.PasswordHash) < 32 {
		return &ErrBadRequest{Message: "fields missing or too short"}
	}
	//email verif pain here!
	emailPrefix, emailPostfix, ok := strings.Cut(user.Email, "@")
	if !ok || emailPrefix == "" || emailPostfix == "" || !strings.Contains(emailPostfix, ".") ||
		strings.Contains(emailPostfix, "@") {
		return &ErrBadRequest{Message: "invalid email address"}
	}

	//TODO: ping email provider and send verification email here

	err := s.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil

}
