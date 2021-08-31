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

func (s *UserService) Create(ctx context.Context, u domain.User) error {
	// Validate user struct before create
	if err := u.Validate(); err != nil {
		return errors.Wrap(err, "UserService - Create - u.Validate")
	}

	// Hash password
	if err := u.HashPassword(); err != nil {
		return errors.Wrap(err, "UserService - Create - u.HashPassword")
	}

	// Create a new user in database
	if err := s.repo.Create(ctx, u); err != nil {
		return errors.Wrap(err, "UserService - Create - s.repo.Create")
	}

	return nil
}
