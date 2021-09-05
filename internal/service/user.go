package service

import (
	"context"
	"github.com/pkg/errors"

	"github.com/ysomad/go-auth-service/internal/domain"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{r}
}

func (s *UserService) Create(ctx context.Context, u *domain.User) error {
	// Encrypt password
	if err := u.EncryptPassword(); err != nil {
		return err
	}

	// Create a new user in database
	if err := s.repo.Create(ctx, u); err != nil {
		return err
	}

	u.Sanitize()

	return nil
}

func (s *UserService) authenticate(ctx context.Context, u *domain.User) error {
	// Get encrypted user password
	encryptedPwd, err := s.repo.GetPassword(ctx, u.ID)
	if err != nil {
		return err
	}

	// Compare password with encrypted password
	if !u.CompareHashAndPassword(encryptedPwd) {
		return errors.New("incorrect password")
	}

	return nil
}

// UpdateState updates User is_active flag
func (s *UserService) UpdateState(ctx context.Context, u *domain.User) error {
	if err := s.repo.UpdateState(ctx, u); err != nil {
		return err
	}

	return nil
}

// Update updates User field values with new values if password is correct
func (s *UserService) Update(ctx context.Context, u *domain.User) error {
	if err := s.authenticate(ctx, u); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return err
	}

	u.Sanitize()

	return nil
}
