package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	User interface {
		SignUp(context.Context, *domain.CreateUserRequest) (*domain.User, error)
		Archive(context.Context, *domain.ArchiveUserRequest) error // TODO: return domain.User
		Update(context.Context, *domain.UpdateUserRequest) (*domain.User, error)
		GetByID(context.Context, int) (*domain.User, error)
	}

	UserRepo interface {
		Create(context.Context, *domain.User) error
		Archive(context.Context, *domain.ArchiveUserRequest) error // TODO: work with domain.User model
		Update(context.Context, *domain.User) error
		GetByID(context.Context, *domain.User) error
	}
)
