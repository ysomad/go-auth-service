package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	User interface {
		SignUp(context.Context, *domain.CreateUserRequest) (*domain.CreateUserResponse, error)
		UpdateState(context.Context, *domain.UpdateStateUserRequest) (*domain.UpdateStateUserResponse, error) // TODO: implement UpdateState service
		Update(context.Context, *domain.User) error
	}

	UserRepo interface {
		Insert(context.Context, *domain.CreateUserRequest) (*domain.CreateUserResponse, error)
		GetPassword(context.Context, int) (string, error)
		UpdateState(context.Context, *domain.UpdateStateUserResponse) error // TODO: implement UpdateState reposotiry
		Update(context.Context, *domain.User) error
	}
)
