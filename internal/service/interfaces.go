package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/ysomad/go-auth-service/internal/entity"
)

type (
	User interface {
		Create(ctx context.Context, email string, password string) error
		Archive(ctx context.Context, id uuid.UUID, isArchive bool) error
		PartialUpdate(ctx context.Context, u entity.UserPartialUpdateDTO) error
		GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	}

	UserRepo interface {
		Create(ctx context.Context, email string, password string) error
		Archive(ctx context.Context, id uuid.UUID, isArchive bool) error
		PartialUpdate(ctx context.Context, id uuid.UUID, cols entity.UpdateColumns) error
		GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
	}

	Auth interface {
		Login(ctx context.Context, cred entity.UserCredentialsDTO, security entity.SessionSecurityDTO) (entity.JWT, error)
		RefreshToken(ctx context.Context, security entity.SessionSecurityDTO) (entity.JWT, error)
	}

	SessionRepo interface {
		Create(ctx context.Context, s *entity.Session) error
		GetOne(ctx context.Context, refreshToken uuid.UUID) (*entity.Session, error)
		Terminate(ctx context.Context, refreshToken uuid.UUID) error
	}
)
