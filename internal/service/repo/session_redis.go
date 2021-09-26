package repo

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/ysomad/go-auth-service/internal/entity"
)

type sessionRepo struct {
	*redis.Client
}

func NewSessionRepo(r *redis.Client) *sessionRepo {
	return &sessionRepo{r}
}

func (r *sessionRepo) sessionKey(userID uuid.UUID, refreshToken uuid.UUID) string {
	return fmt.Sprintf("%s:%s", userID, refreshToken)
}

// sessionList returns list of found session with refresh token
func (r *sessionRepo) sessionList(ctx context.Context, refreshToken uuid.UUID, cursor uint64, count int64) ([]string, uint64, error) {
	return r.Scan(ctx, cursor, fmt.Sprintf("*:%s", refreshToken), count).Result()
}

// Create sets new refresh session to redis with refresh token as key
func (r *sessionRepo) Create(ctx context.Context, s *entity.Session) error {
	b, err := s.MarshalBinary()
	if err != nil {
		return err
	}

	if err = r.Set(ctx, r.sessionKey(s.UserID, s.RefreshToken), b, s.ExpiresIn).Err(); err != nil {
		return err
	}

	return nil
}

func (r *sessionRepo) GetOne(ctx context.Context, refreshToken uuid.UUID) (*entity.Session, error) {
	sessionKeys, _, err := r.sessionList(ctx, refreshToken, 0, 0)
	if err != nil {
		return nil, err
	}

	if len(sessionKeys) == 0 {
		return nil, entity.ErrSessionExpired
	}

	var session entity.Session

	if err = r.Get(ctx, sessionKeys[0]).Scan(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepo) Terminate(ctx context.Context, refreshToken uuid.UUID) error {
	sessionKeys, _, err := r.sessionList(ctx, refreshToken, 0, 0)
	if err != nil {
		return err
	}

	if len(sessionKeys) == 0 {
		return entity.ErrSessionExpired
	}

	if err = r.Del(ctx, sessionKeys[0]).Err(); err != nil {
		return err
	}

	return nil
}
