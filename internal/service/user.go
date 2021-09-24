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

// SignUp creates new user with email and encrypted password
func (s *UserService) SignUp(ctx context.Context, req entity.CreateUserRequest) error {
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 11)
	if err != nil {
		return err
	}

	if err = s.repo.Create(ctx, req.Email, string(password)); err != nil {
		return err
	}

	return nil
}

// Archive sets user is_archive
func (s *UserService) Archive(ctx context.Context, id uuid.UUID, isArchive bool) error {
	if err := s.repo.Archive(ctx, id, isArchive); err != nil {
		return err
	}

	return nil
}

// PartialUpdate updates all updatable user columns
func (s *UserService) PartialUpdate(ctx context.Context, id uuid.UUID, req entity.PartialUpdateRequest) error {
	cols := map[string]interface{}{
		"username":   req.Username,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
	}

	for k, v := range cols {
		if v == "" || v == nil {
			delete(cols, k)
		}
	}

	if len(cols) == 0 {
		return entity.PartialUpdateErr
	}

	if err := s.repo.PartialUpdate(ctx, id, cols); err != nil {
		return err
	}

	return nil
}

// GetByID gets user data by ID
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
