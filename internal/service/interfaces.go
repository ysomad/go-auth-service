package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/ysomad/go-auth-service/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=service_test

type (
	User interface {
		SignUp(ctx context.Context, req entity.CreateUserRequest) error
		Archive(ctx context.Context, id uuid.UUID, isArchive bool) error
		PartialUpdate(ctx context.Context, id uuid.UUID, req entity.PartialUpdateRequest) error
		GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	}

	UserRepo interface {
		Create(ctx context.Context, email string, password string) error
		Archive(ctx context.Context, id uuid.UUID, isArchive bool) error
		PartialUpdate(ctx context.Context, id uuid.UUID, cols map[string]interface{}) error
		GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
	}

	Auth interface {
		Login(ctx context.Context, req entity.LoginRequest, dto entity.SessionSecurityDTO) (entity.LoginResponse, error)
		RefreshToken(ctx context.Context, dto entity.SessionSecurityDTO) (entity.LoginResponse, error)
	}

	SessionRepo interface {
		Create(ctx context.Context, s entity.Session) error
		Get(ctx context.Context, refreshToken uuid.UUID) (entity.Session, error)
		Terminate(ctx context.Context, refreshToken uuid.UUID) error
	}
)
