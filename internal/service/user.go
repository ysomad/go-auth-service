package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/entity"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{r}
}

// Create creates new user with email and encrypted password
func (s *UserService) Create(ctx context.Context, email string, password string) (*entity.User, error) {
	p, err := entity.EncryptPassword(password)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.Create(ctx, email, p)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Archive updates is_archive field
func (s *UserService) Archive(ctx context.Context, id int, isArchive bool) error {
	if err := s.repo.Archive(ctx, id, isArchive); err != nil {
		return err
	}

	return nil
}

// PartialUpdate updates all updatable user columns
func (s *UserService) PartialUpdate(ctx context.Context, id int, cols map[string]interface{}) (*entity.User, error) {
	u, err := s.repo.PartialUpdate(ctx, id, cols)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetByID gets user data by ID
func (s *UserService) GetByID(ctx context.Context, id int) (*entity.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
