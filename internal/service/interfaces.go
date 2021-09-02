package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	User interface {
		Create(context.Context, *domain.User) error
		Archive(context.Context, *domain.User) error
	}

	UserRepo interface {
		Create(context.Context, *domain.User) error
		GetPassword(context.Context, int) (string, error)
		Archive(context.Context, int) error
	}
)
