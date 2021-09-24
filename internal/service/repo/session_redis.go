package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ysomad/go-auth-service/internal/entity"
)

type SessionRepo struct {
	*redis.Client
}

func NewSessionRepo(r *redis.Client) *SessionRepo {
	return &SessionRepo{r}
}

// Create sets new refresh session to redis with refresh token as key
func (r *SessionRepo) Create(ctx context.Context, s entity.Session) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", s.UserID.String(), s.RefreshToken.String())

	// Create session
	if err = r.SAdd(ctx, key, b).Err(); err != nil {
		return err
	}

	// Set expiry
	if err = r.Expire(ctx, key, s.ExpiresIn).Err(); err != nil {
		return err
	}

	return nil
}
