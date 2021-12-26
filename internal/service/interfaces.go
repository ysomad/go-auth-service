package service

import (
	"context"

	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	Account interface {
		// Create new account with email and password,
		Create(ctx context.Context, email, password string) error

		// GetByID account.
		GetByID(ctx context.Context, aid string) (domain.Account, error)

		// GetByEmail account.
		GetByEmail(ctx context.Context, email string) (domain.Account, error)

		// Archive sets account IsArchive state to true.
		Archive(ctx context.Context, aid string) error
	}

	AccountRepo interface {
		// Create account with given credentials.
		Create(ctx context.Context, a domain.Account) error

		// FindByID account in DB.
		FindByID(ctx context.Context, aid string) (domain.Account, error)

		// FindByEmail account in DB.
		FindByEmail(ctx context.Context, email string) (domain.Account, error)

		// Archive sets entity.Account.IsArchive state to provided value.
		Archive(ctx context.Context, aid string, archive bool) error
	}

	Auth interface {
		// EmailLogin creates new session using provided account email and password.
		EmailLogin(ctx context.Context, email, password string, d domain.Device) (domain.SessionCookie, error)

		// Logout logs out session by id.
		Logout(ctx context.Context, sid string) error

		// GetAccessToken generates JWT token which must be used to complete protected operations.
		NewAccessToken(ctx context.Context, aid, password string) (domain.Token, error)

		// ParseAccessToken parses and validates JWT access token, returns subject from payload.
		ParseAccessToken(ctx context.Context, token string) (string, error)
	}

	Session interface {
		// Create new session for account with id and device.
		Create(ctx context.Context, aid string, d domain.Device) (domain.Session, error)

		// Get session by id.
		Get(ctx context.Context, sid string) (domain.Session, error)

		// GetAll account sessions using provided account id.
		GetAll(ctx context.Context, aid string) ([]domain.Session, error)

		// Terminate sessions by id.
		Terminate(ctx context.Context, sid string) error

		// TerminateAll account sessions excluding current one.
		TerminateAll(ctx context.Context, aid, currSid string) error
	}

	SessionRepo interface {
		// Create new session in DB.
		Create(ctx context.Context, s domain.Session) error

		// Get session by id.
		Get(ctx context.Context, sid string) (domain.Session, error)

		// GetAll accounts sessions by provided account id.
		GetAll(ctx context.Context, aid string) ([]domain.Session, error)

		// Delete session by id.
		Delete(ctx context.Context, sid string) error

		// DeleteAll account sessions by provided account id excluding current session.
		DeleteAll(ctx context.Context, aid, currSid string) error
	}
)
