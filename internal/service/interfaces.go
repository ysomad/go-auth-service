package service

import (
	"context"
	"time"

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
		Create(ctx context.Context, cred domain.AccountCredentials) error

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
		GetAccessToken(ctx context.Context, aid string) (domain.Token, error)
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

		// TerminateAll account sessions using provided account id.
		TerminateAll(ctx context.Context, aid string) error
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

		// DeleteAll account sessions by provided account id.
		DeleteAll(ctx context.Context, aid string) error
	}

	CacheRepo interface {
		// Set the given key/value in cache,
		// overwriting any existing value associated with that key.
		Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error

		// Add the given key/value to cache ONLY IF the key does not already exist.
		Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error

		// Get content associated with the given key from cache.
		// Decoding it into the given pointer.
		Get(ctx context.Context, key string, pointer interface{}) error

		// Delete the given key from cache.
		Delete(ctx context.Context, key string) error
	}
)
