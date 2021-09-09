package service

import (
	"context"

	"github.com/ysomad/go-auth-service/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=service_test

type (
	User interface {
		SignUp(ctx context.Context, req entity.CreateUserRequest) error
		Archive(ctx context.Context, id int, isArchive bool) error
		PartialUpdate(ctx context.Context, id int, req entity.PartialUpdateRequest) error
		GetByID(ctx context.Context, id int) (*entity.User, error)
	}

	UserRepo interface {
		Create(ctx context.Context, email string, password string) error
		Archive(ctx context.Context, id int, isArchive bool) error
		PartialUpdate(ctx context.Context, id int, cols map[string]interface{}) error
		GetByID(ctx context.Context, id int) (*entity.User, error)
	}
)
