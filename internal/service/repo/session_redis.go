package repo

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/ysomad/go-auth-service/internal/entity"
)

type SessionRepo struct {
	*redis.Client
}

func NewSessionRepo(r *redis.Client) *SessionRepo {
	return &SessionRepo{r}
}

// Create sets new refresh session to redis with refresh token as key
func (r *SessionRepo) Create(s entity.RefreshSession) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err = r.Set(s.RefreshToken.String(), b, s.ExpiresIn).Err(); err != nil {
		return err
	}

	return nil
}
