// Package service implements application business logic. Each logic group in own file.
package service

import (
	"context"

	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	Translation interface {
		Translate(context.Context, domain.Translation) (domain.Translation, error)
		History(context.Context) ([]domain.Translation, error)
	}

	TranslationRepo interface {
		Store(context.Context, domain.Translation) error
		GetHistory(context.Context) ([]domain.Translation, error)
	}

	TranslationWebAPI interface {
		Translate(domain.Translation) (domain.Translation, error)
	}

	// User interfaces

	User interface {
		Create(context.Context, domain.User) error
	}

	UserRepo interface {
		Create(context.Context, domain.User) error
	}

	JWTAuthentication interface {
	}
)
