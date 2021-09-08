package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/entity"
)

type (
	User interface {
		SignUp(context.Context, *entity.CreateUserRequest) (*entity.User, error)
		Archive(context.Context, *entity.ArchiveUserRequest) error // TODO: return entity.User
		Update(context.Context, *entity.UpdateUserRequest) (*entity.User, error)
		GetByID(context.Context, int) (*entity.User, error)
	}

	UserRepo interface {
		Create(context.Context, *entity.User) error
		Archive(context.Context, *entity.ArchiveUserRequest) error // TODO: work with entity.User model
		Update(context.Context, *entity.User) error
		GetByID(context.Context, *entity.User) error
	}
)
