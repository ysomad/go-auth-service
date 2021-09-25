package service

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ysomad/go-auth-service/internal/entity"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{r}
}

// Create creates new userRepo with email and encrypted password
func (s *UserService) Create(ctx context.Context, email string, password string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return err
	}

	if err = s.repo.Create(ctx, email, string(b)); err != nil {
		return err
	}

	return nil
}

// Archive sets userRepo is_archive
func (s *UserService) Archive(ctx context.Context, id uuid.UUID, isArchive bool) error {
	if err := s.repo.Archive(ctx, id, isArchive); err != nil {
		return err
	}

	return nil
}

// PartialUpdate updates all updatable userRepo columns
func (s *UserService) PartialUpdate(ctx context.Context, u entity.UserPartialUpdateDTO) error {
	// Validate update columns
	cols := entity.UpdateColumns{
		"username":   u.Username,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
	}
	if err := cols.Validate(); err != nil {
		return err
	}

	if err := s.repo.PartialUpdate(ctx, u.ID, cols); err != nil {
		return err
	}

	return nil
}

// GetByID gets userRepo data by ID
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
