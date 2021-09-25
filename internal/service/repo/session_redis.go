package repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/ysomad/go-auth-service/internal/entity"
)

const (
	sessionPrefix = "sess"
)

type SessionRepo struct {
	*redis.Client
}

func NewSessionRepo(r *redis.Client) *SessionRepo {
	return &SessionRepo{r}
}

func (r *SessionRepo) sessionKey(key string) string {
	return fmt.Sprintf("%s:%s", sessionPrefix, key)
}

// Create sets new refresh session to redis with refresh token as key
func (r *SessionRepo) Create(ctx context.Context, s entity.Session) error {
	tokenString := s.RefreshToken.String()
	key := r.sessionKey(tokenString)

	if err := r.HSet(ctx, key, map[string]interface{}{
		"token":   tokenString,
		"uid":     s.UserID.String(),
		"ua":      s.UserAgent,
		"ip":      s.UserIP,
		"fp":      s.Fingerprint.String(),
		"exp":     s.ExpiresAt,
		"created": s.CreatedAt,
	}).Err(); err != nil {
		return err
	}

	if err := r.Expire(ctx, key, s.ExpiresIn).Err(); err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) Get(ctx context.Context, refreshToken uuid.UUID) (entity.Session, error) {
	res, err := r.HGetAll(ctx, r.sessionKey(refreshToken.String())).Result()
	if err != nil {
		return entity.Session{}, err
	}

	if res["token"] == "" {
		return entity.Session{}, entity.ErrSessionExpired
	}

	// Parse values from strings
	token, err := uuid.Parse(res["token"])
	if err != nil {
		return entity.Session{}, err
	}

	uid, err := uuid.Parse(res["uid"])
	if err != nil {
		return entity.Session{}, err
	}

	fp, err := uuid.Parse(res["fp"])
	if err != nil {
		return entity.Session{}, err
	}

	exp, err := strconv.ParseInt(res["exp"], 10, 64)
	if err != nil {
		return entity.Session{}, err
	}

	created, err := time.Parse(time.RFC3339Nano, res["created"])
	if err != nil {
		return entity.Session{}, err
	}

	s := entity.Session{
		RefreshToken: token,
		UserID: uid,
		UserAgent: res["ua"],
		UserIP: res["ip"],
		Fingerprint: fp,
		ExpiresAt: exp,
		CreatedAt: created,
	}

	return s, nil
}

func (r *SessionRepo) Terminate(ctx context.Context, refreshToken uuid.UUID) error {
	if err := r.Del(ctx, r.sessionKey(refreshToken.String())).Err(); err != nil {
		return err
	}

	return nil
}
