package service

import (
	"context"
	"net/url"

	"github.com/ysomad/go-auth-service/internal/domain"
)

type (
	Account interface {
		// Create new account, username, email and password should be provided, returns account id.
		Create(ctx context.Context, a domain.Account) (string, error)

		// GetByID account.
		GetByID(ctx context.Context, aid string) (domain.Account, error)

		// GetByEmail account.
		GetByEmail(ctx context.Context, email string) (domain.Account, error)

		// Delete sets account IsArchive state to true.
		Delete(ctx context.Context, aid string) error

		// Verify verifies account using provided code.
		Verify(ctx context.Context, code string) error
	}

	AccountRepo interface {
		// Create account with given credentials, returns id of created account.
		Create(ctx context.Context, a domain.Account) (string, error)

		// FindByID account in DB.
		FindByID(ctx context.Context, aid string) (domain.Account, error)

		// FindByEmail account in DB.
		FindByEmail(ctx context.Context, email string) (domain.Account, error)

		// Archive sets entity.Account.IsArchive state to provided value.
		Archive(ctx context.Context, aid string, archive bool) error
	}

	Auth interface {
		// EmailLogin creates new session using provided account email and password.
		EmailLogin(ctx context.Context, email, password string, d Device) (SessionCookie, error)

		// Logout logs out session by id.
		Logout(ctx context.Context, sid string) error

		// NewAccessToken generates JWT token which must be used to perform protected operations.
		NewAccessToken(ctx context.Context, aid, password string) (string, error)

		// ParseAccessToken parses and validates JWT access token, returns subject from payload.
		ParseAccessToken(ctx context.Context, t string) (string, error)
	}

	SocialAuth interface {
		// AuthorizationURL returns OAuth authorization URL of given provider with
		// client id, scope and state query parameters.
		AuthorizationURL(ctx context.Context, provider string) (*url.URL, error)

		// GitHubLogin handles OAuth2 login via GitHub.
		GitHubLogin(ctx context.Context, code string, d Device) (SessionCookie, error)

		// GoogleLogin handles OAuth2 login via Google.
		GoogleLogin(ctx context.Context, code string, d Device) (SessionCookie, error)
	}

	Session interface {
		// Create new session for account with id and device of given provider.
		Create(ctx context.Context, aid, provider string, d Device) (domain.Session, error)

		// GetByID session.
		GetByID(ctx context.Context, sid string) (domain.Session, error)

		// GetAll account sessions using provided account id.
		GetAll(ctx context.Context, aid string) ([]domain.Session, error)

		// Terminate session by id excluding current one.
		Terminate(ctx context.Context, sid, currSid string) error

		// TerminateAll account sessions excluding current one.
		TerminateAll(ctx context.Context, aid, sid string) error
	}

	SessionRepo interface {
		// Create new session in DB.
		Create(ctx context.Context, s domain.Session) error

		// FindByID session.
		FindByID(ctx context.Context, sid string) (domain.Session, error)

		// FindAll accounts sessions by provided account id.
		FindAll(ctx context.Context, aid string) ([]domain.Session, error)

		// Delete session by id.
		Delete(ctx context.Context, sid string) error

		// DeleteAll account sessions by provided account id excluding current session.
		DeleteAll(ctx context.Context, aid, sid string) error
	}
)
