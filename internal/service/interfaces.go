package service

import (
	"context"
	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	User interface {
		Create(context.Context, *domain.User) (error, map[string]string)
	}

	UserRepo interface {
		Create(context.Context, *domain.User) error
	}

	JWTAuthentication interface {
	}
)
