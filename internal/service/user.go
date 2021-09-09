package service

import (
	"context"
	"errors"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/pkg/crypto"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{r}
}

// Create creates new user with email and encrypted password
func (s *UserService) Create(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error) {
	p, err := crypto.EncryptPassword(req.Password, 11)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.Create(ctx, req.Email, p)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Archive sets user is_archive
func (s *UserService) Archive(ctx context.Context, id int, isArchive bool) error {
	if err := s.repo.Archive(ctx, id, isArchive); err != nil {
		return err
	}

	return nil
}

// PartialUpdate updates all updatable user columns
func (s *UserService) PartialUpdate(ctx context.Context, id int, req entity.PartialUpdateRequest) (*entity.User, error) {
	cols := map[string]interface{}{
		"username":   req.Username,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
	}

	for k, v := range cols {
		if v == nil || v == "" {
			delete(cols, k)
		}
	}

	if len(cols) == 0 {
		return nil, errors.New("provide at least one field to update resource partially")
	}

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
