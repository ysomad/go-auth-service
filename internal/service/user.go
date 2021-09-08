package service

import (
	"context"
	"time"

	"github.com/ysomad/go-auth-service/internal/entity"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{r}
}

func (s *UserService) Create(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error) {
	p, err := entity.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.Create(ctx, req.Email, p)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Archive updates is_archive field for user
func (s *UserService) Archive(ctx context.Context, req *entity.ArchiveUserRequest) error {
	if err := s.repo.Archive(ctx, req); err != nil {
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

	if err := s.repo.Update(ctx, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

// GetByID gets user data by ID
func (s *UserService) GetByID(ctx context.Context, id int) (*entity.User, error) {
	u := entity.User{ID: id}

	if err := s.repo.GetByID(ctx, &u); err != nil {
		return nil, err
	}

	return &u, nil
}
