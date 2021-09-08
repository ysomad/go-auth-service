package service

import (
	"context"
	"time"

	"github.com/ysomad/go-auth-service/internal/entity"
)

type UserService struct {
	r UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{r}
}

func (s *UserService) SignUp(ctx context.Context, req *entity.CreateUserRequest) (*entity.User, error) {
	u := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := u.EncryptPassword(); err != nil {
		return nil, err
	}

	if err := s.r.Create(ctx, &u); err != nil {
		return nil, err
	}

	u.Sanitize()

	return &u, nil
}

// Archive updates is_archive field for user
func (s *UserService) Archive(ctx context.Context, req *entity.ArchiveUserRequest) error {
	if err := s.r.Archive(ctx, req); err != nil {
		return err
	}

	return nil
}

// Update updates all updatable user columns
func (s *UserService) Update(ctx context.Context, req *entity.UpdateUserRequest) (*entity.User, error) {
	u := entity.User{
		ID:        req.ID,
		Username:  &req.Username,
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		UpdatedAt: time.Now(),
	}

	if err := s.r.Update(ctx, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

// GetByID gets user data by ID
func (s *UserService) GetByID(ctx context.Context, id int) (*entity.User, error) {
	u := entity.User{ID: id}

	if err := s.r.GetByID(ctx, &u); err != nil {
		return nil, err
	}

	return &u, nil
}
