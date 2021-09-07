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

func (s *UserService) SignUp(ctx context.Context, u *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	if err := u.EncryptPassword(); err != nil {
		return nil, err
	}

	resp, err := s.repo.Insert(ctx, u)
	if err != nil {
		return nil, err
	}

	u.Sanitize()

	return resp, nil
}

// UpdateState updates User is_active flag
func (s *UserService) UpdateState(ctx context.Context, u *domain.UpdateStateUserRequest) (*domain.UpdateStateUserResponse, error) {
	resp := domain.UpdateStateUserResponse{
		ID: u.ID,
		IsActive: *u.IsActive,
		UpdatedAt: time.Now(),
	}

	if err := s.repo.UpdateState(ctx, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update updates User field values with new values if password is correct
func (s *UserService) Update(ctx context.Context, u *domain.User) error {
	if err := s.repo.Update(ctx, u); err != nil {
		return err
	}

	return nil
}
