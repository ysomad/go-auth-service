package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	User interface {
		Create(context.Context, *domain.User) error
		UpdateState(context.Context, *domain.User) error
		Update(context.Context, *domain.User) error
	}

	UserRepo interface {
		Create(context.Context, *domain.User) error
		GetPassword(context.Context, int) (string, error)
		UpdateState(context.Context, *domain.User) error
		Update(context.Context, *domain.User) error
	}
)
