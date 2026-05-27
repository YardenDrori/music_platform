package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/identity"
)

type service struct {
	repo repository
}

func NewService(repo repository) Service {
	return &service{
		repo: repo,
	}
}

func requireSelf(ctx context.Context, id uuid.UUID) error {
	requesterID, ok := identity.UserIDFromContext(ctx)
	if !ok {
		return ErrUnathenticated
	}
	if id != requesterID {
		return ErrForbidden
	}
	return nil
}

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

func (s *service) FindByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, &ErrBadRequest{Message: "email not provided"}
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err = requireSelf(ctx, user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) FindByUsername(ctx context.Context, username string) (*User, error) {
	if username == "" {
		return nil, &ErrBadRequest{Message: "username not provided"}
	}

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err = requireSelf(ctx, user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	if err := requireSelf(ctx, id); err != nil {
		return err
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeactivateAccount(ctx context.Context, id uuid.UUID) error {
	if err := requireSelf(ctx, id); err != nil {
		return err
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return err
		}
		return fmt.Errorf("fetching account info for deactivation: %w", err)
	}

	user.Active = false

	err = s.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("deactivating account: %w", err)
	}

	return nil
}
