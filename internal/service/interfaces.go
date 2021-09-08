package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=service_test

type (
	User interface {
		Create(ctx context.Context, email string, password string) (*entity.User, error)
		Archive(ctx context.Context, id int, isArchive bool) error
		Update(context.Context, *entity.UpdateUserRequest) (*entity.User, error)
		GetByID(context.Context, int) (*entity.User, error)
	}

	UserRepo interface {
		Create(ctx context.Context, email string, password string) (*entity.User, error)
		Archive(ctx context.Context, id int, isArchive bool) error
		Update(context.Context, *entity.User) error
		GetByID(context.Context, *entity.User) error
	}
)
