package service

import (
	"context"
	"time"

	"github.com/ysomad/go-auth-service/internal/entity"
)

type (
	User interface {
		// Register new user with email and password.
		Register(ctx context.Context, email, password string) (entity.User, error)

		FindByID(ctx context.Context, uid string) (entity.User, error)

		// Archive or restore user.
		Archive(ctx context.Context, uid string, archive bool) error
	}

	Session interface {
		// LoginWithEmail creates new user session.
		LoginWithEmail(ctx context.Context, email, password string,
			d entity.Device) (entity.Session, error)

		// Find session by id.
		Find(ctx context.Context, sid string) (entity.Session, error)

		// FindAll sessions by user id.
		FindAll(ctx context.Context, uid string) ([]entity.Session, error)

		// Terminate sessions by id.
		Terminate(ctx context.Context, sid string) error

		// TerminateAll sessions of user with id.
		TerminateAll(ctx context.Context, uid string) error
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

	UserRepo interface {
		Create(ctx context.Context, dto entity.UserSensitiveData) (entity.User, error)
		GetByID(ctx context.Context, uid string) (entity.User, error)
		GetByEmail(ctx context.Context, email string) (entity.User, error)
		Archive(ctx context.Context, uid string, archive bool) (entity.User, error)
	}

	SessionRepo interface {
		// Create new session.
		Create(ctx context.Context, s entity.Session) error

		// Get session by id.
		Get(ctx context.Context, sid string) (entity.Session, error)

		// GetAll sessions by user id.
		GetAll(ctx context.Context, uid string) ([]entity.Session, error)

		// Delete session by id.
		Delete(ctx context.Context, sid string) error

		// DeleteAll sessions by user id.
		DeleteAll(ctx context.Context, uid string) error
	}
)
