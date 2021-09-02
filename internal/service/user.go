package service

import (
	"context"
	"github.com/pkg/errors"

	"github.com/ysomad/go-auth-service/internal/domain"
)

type UserService struct {
	repo     UserRepo
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

// Archive sets user `is_active` column to `false`
func (s *UserService) Archive(ctx context.Context, u *domain.User) error {
	encryptedPwd, err := s.repo.GetPassword(ctx, u.ID)
	if err != nil {
		return err
	}

	// Compare password with encrypted password
	u.SetEncryptedPassword(encryptedPwd)
	if !u.CompareHashAndPassword() {
		return errors.New("incorrect password")
	}

	// Archive user
	if err = s.repo.Archive(ctx, u.ID); err != nil {
		return err
	}

	return nil
}
