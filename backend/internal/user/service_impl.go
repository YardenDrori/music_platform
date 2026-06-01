package user

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	"github.com/YardenDrori/music-platform/internal/identity"
)

type service struct {
	repo           repository
	passwordHasher passwordHasher
}

func NewService(repo repository, hasher passwordHasher) Service {
	return &service{
		repo:           repo,
		passwordHasher: hasher,
	}
}

// errors:
// [ErrUnauthenticated]
// [ErrForbidden]
func requireSelf(ctx context.Context, id uuid.UUID) error {
	requesterID, ok := identity.UserIDFromContext(ctx)
	if !ok {
		return apperrors.ErrUnauthenticated
	}
	if id != requesterID {
		return apperrors.ErrForbidden
	}
	return nil
}

// errors:
// ErrBadRequest
func validateAccountBusinessRules(user *User) error {
	//validation logic
	if utf8.RuneCountInString(user.Username) < 3 || utf8.RuneCountInString(user.Email) < 3 ||
		utf8.RuneCountInString(user.FirstName) < 3 ||
		utf8.RuneCountInString(user.LastName) < 3 ||
		utf8.RuneCountInString(user.PasswordHash) < 32 {
		return &apperrors.ErrBadRequest{Message: "fields missing or too short"}
	}
	//user verif pain here
	if strings.Contains(user.Username, "@") {
		return &apperrors.ErrBadRequest{Message: "invalid username"}
	}

	//email verif pain here!
	emailPrefix, emailPostfix, ok := strings.Cut(user.Email, "@")
	if !ok || emailPrefix == "" || emailPostfix == "" || !strings.Contains(emailPostfix, ".") ||
		strings.Contains(emailPostfix, "@") {
		return &apperrors.ErrBadRequest{Message: "invalid email address"}
	}
	return nil
}

func (s *service) NewAccount(ctx context.Context, user *NewUserRequest) (*User, error) {
	passHash := s.passwordHasher.hashPassword(user.Password)
	newUser := &User{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PasswordHash: passHash,
		CreatedAt:    time.Now().UTC(),
		LastUpdated:  time.Now().UTC(),
		Active:       true,
	}

	if err := validateAccountBusinessRules(newUser); err != nil {
		return nil, err
	}

	//TODO: ping email provider and send verification email here

	err := s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("creating user account: %w", err)
	}

	return newUser, nil
}

func (s *service) Authenticate(
	ctx context.Context,
	identifier string,
	password string,
) (*User, error) {
	var user *User
	var err error
	if strings.Contains(identifier, "@") {
		user, err = s.findByEmail(ctx, identifier)
	} else {
		user, err = s.findByUsername(ctx, identifier)
	}
	if err != nil {
		return nil, fmt.Errorf("authenticating user: %w", apperrors.ErrNotFound)
	}

	passwordsMatch, err := s.passwordHasher.verifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticating user: %w", apperrors.ErrUnauthenticated)
	}
	if passwordsMatch {
		return user, nil
	}
	return nil, apperrors.ErrUnauthenticated
}

func (s *service) findByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, &apperrors.ErrBadRequest{Message: "email not provided"}
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("finding user by email: %w", err)
	}

	return user, nil
}

func (s *service) findByUsername(ctx context.Context, username string) (*User, error) {
	if username == "" {
		return nil, &apperrors.ErrBadRequest{Message: "username not provided"}
	}

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("finding user by username: %w", err)
	}

	return user, nil
}

// errors:
// [ErrUnauthenticated]
// [ErrForbidden]
// [ErrNotFound]
// [errorf]
func (s *service) FindByUUID(ctx context.Context, id uuid.UUID) (*User, error) {
	if err := requireSelf(ctx, id); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding user by id: %w", err)
	}

	return user, nil
}

// errors:
// ErrBadRequest,
// ErrConflict
// ErrForbidden
// ErrUnathenticated
// errorf
func (s *service) UpdateAccount(ctx context.Context, user *NewUserRequest) error {
	if err := requireSelf(ctx, user.ID); err != nil {
		return err
	}

	currUserInfo, err := s.FindByUUID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("updating user profile: fetching account info: %w", err)
	}

	passwordsMatch, err := s.passwordHasher.verifyPassword(user.Password, currUserInfo.PasswordHash)
	if err != nil {
		return fmt.Errorf("updating user profile: checking if passwords match: %w", err)
	}

	var newPass string
	if passwordsMatch {
		newPass = currUserInfo.PasswordHash
	} else {
		newPass = s.passwordHasher.hashPassword(user.Password)
	}

	updatedUser := &User{
		ID:           user.ID,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PasswordHash: newPass,
		CreatedAt:    currUserInfo.CreatedAt,
		Active:       currUserInfo.Active,
	}

	if err := validateAccountBusinessRules(updatedUser); err != nil {
		return err
	}

	err = s.repo.Update(ctx, updatedUser)
	if err != nil {
		return fmt.Errorf("updating user account: %w", err)
	}

	return nil
}

func (s *service) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	if err := requireSelf(ctx, id); err != nil {
		return err
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting user account: %w", err)
	}

	return nil
}

// errors:
// ErrUnathenticated,
// ErrForbidden,
// ErrConflict,
// errorf
func (s *service) DeactivateAccount(ctx context.Context, id uuid.UUID) error {
	if err := requireSelf(ctx, id); err != nil {
		return err
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("fetching account for deactivation: %w", err)
	}

	user.Active = false

	err = s.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("deactivating user account: %w", err)
	}

	return nil
}
