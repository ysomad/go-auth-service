package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	User interface {
		SignUp(context.Context, *domain.CreateUserRequest) (*domain.CreateUserResponse, error)
		Archive(context.Context, *domain.ArchiveUserRequest) (*domain.ArchiveUserResponse, error)
		PartialUpdate(context.Context) error // TODO: implement PartialUpdate
	}

	UserRepo interface {
		Insert(context.Context, *domain.CreateUserRequest) (*domain.CreateUserResponse, error)
		GetPassword(context.Context, int) (string, error)
		Archive(context.Context, *domain.ArchiveUserResponse) error
		PartialUpdate(context.Context) error // TODO: implement PartialUpdate
	}
)
