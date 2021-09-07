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

// Archive updates is_archive field for user
func (s *UserService) Archive(ctx context.Context, req *domain.ArchiveUserRequest) (*domain.ArchiveUserResponse, error) {
	resp := domain.ArchiveUserResponse{
		ID: req.ID,
		IsArchive: *req.IsArchive,
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Archive(ctx, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *UserService) PartialUpdate(ctx context.Context) error {
	return nil
}
