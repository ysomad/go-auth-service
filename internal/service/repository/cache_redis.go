package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
)

type cacheRepo struct {
	cli *redis.Client
}

func NewCacheRepo(r *redis.Client) *cacheRepo {
	return &cacheRepo{r}
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
	count, err := r.cli.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("r.cli.Exists.Result: %w", err)
	}

	if count == 0 {
		return r.Set(ctx, key, val, ttl)
	}

	return apperrors.ErrCacheDuplicate
}

func (r *cacheRepo) Get(ctx context.Context, key string, pointer interface{}) error {
	b, err := r.cli.Get(ctx, key).Bytes()
	if err != nil {
		return fmt.Errorf("r.cli.Get.Bytes: %w", err)
	}

	if err == redis.Nil {
		return apperrors.ErrCacheNotFound
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
