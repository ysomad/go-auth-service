package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ErrNotFound  = errors.New("given key was not found in cache")
	ErrNotStored = errors.New("given key cannot be cached")
)

type cacheRepo struct {
	cli *redis.Client
}

func NewCacheRepo(r *redis.Client) *cacheRepo {
	return &cacheRepo{r}
}

// BuildCacheKey separated by colon.
func BuildCacheKey(left string, right string) string {
	return fmt.Sprintf("%s:%s", left, right)
}

func (r *cacheRepo) Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if err = r.cli.Set(ctx, key, b, ttl).Err(); err != nil {
		return fmt.Errorf("r.cli.Set.Err: %w", err)
	}

	return nil
}

func (r *cacheRepo) Add(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	exists, err := r.cli.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("r.cli.Exists.Result: %w", err)
	}

	if exists == 0 {
		return r.Set(ctx, key, val, ttl)
	}

	return ErrNotStored
}

func (r *cacheRepo) Get(ctx context.Context, key string, pointer interface{}) error {
	b, err := r.cli.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return fmt.Errorf("r.cli.Get.Bytes: %w", ErrNotFound)
	}

	if err != nil {
		return fmt.Errorf("r.cli.Get.Bytes: %w", err)
	}

	if err = json.Unmarshal(b, pointer); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}

func (r *cacheRepo) Delete(ctx context.Context, key string) error {
	if err := r.cli.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("r.cli.Del.Err: %w", err)
	}

	return nil
}
