package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	User interface {
		SignUp(context.Context, *domain.CreateUserRequest) (*domain.User, error)
		Archive(context.Context, *domain.ArchiveUserRequest) error
		Update(context.Context, *domain.UpdateUserRequest) (*domain.User, error)
	}

	UserRepo interface {
		Create(context.Context, *domain.User) error
		Archive(context.Context, *domain.ArchiveUserRequest) error
		Update(context.Context, *domain.User) error
	}
)
