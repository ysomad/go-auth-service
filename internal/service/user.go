package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/domain"
	"time"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{r}
}

func (s *UserService) SignUp(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	u := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := u.EncryptPassword(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, &u); err != nil {
		return nil, err
	}

	u.Sanitize()

	return &u, nil
}

// Archive updates is_archive field for user
func (s *UserService) Archive(ctx context.Context, req *domain.ArchiveUserRequest) error {
	if err := s.repo.Archive(ctx, req); err != nil {
		return err
	}

	return nil
}

// Update updates all updatable user columns
func (s *UserService) Update(ctx context.Context, req *domain.UpdateUserRequest) (*domain.User, error) {
	u := domain.User{
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
