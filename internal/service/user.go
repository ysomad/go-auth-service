package service

import (
	"context"
	"fmt"
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

func (s *UserService) hash(str string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(str), 11)
}

// Create creates new userRepo with email and encrypted password
func (s *UserService) Create(ctx context.Context, email string, password string) error {
	b, err := s.hash(password)
	if err != nil {
		return fmt.Errorf("UserService - Create - s.hash: %w", err)
	}

	if err = s.repo.Create(ctx, email, string(b)); err != nil {
		return fmt.Errorf("UserService - Create - s.repo.Create: %w", err)
	}

	return nil
}

// Archive sets userRepo is_archive
func (s *UserService) Archive(ctx context.Context, id uuid.UUID, isArchive bool) error {
	if err := s.repo.Archive(ctx, id, isArchive); err != nil {
		return fmt.Errorf("UserService - Archive - s.repo.Archive: %w", err)
	}

	return nil
}

// PartialUpdate updates all updatable userRepo columns
func (s *UserService) PartialUpdate(ctx context.Context, id uuid.UUID, cols entity.UpdateColumns) error {
	if err := s.repo.PartialUpdate(ctx, id, cols); err != nil {
		return fmt.Errorf("UserService - PartialUpdate - s.repo.PartialUpdate: %w", err)
	}

	return nil
}

// GetByID gets userRepo data by ID
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserService - GetByID - s.repo.GetByID: %w", err)
	}

	return u, nil
}
