package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service/repository"
)

const (
	userCacheKey = "user"
)

type userService struct {
	userRepo UserRepo

	cache    CacheRepo
	cacheTTL time.Duration
}

func NewUserService(r UserRepo, c CacheRepo, cacheTTL time.Duration) *userService {
	return &userService{r, c, cacheTTL}
}

func (s *userService) FindByID(ctx context.Context, uid string) (entity.User, error) {
	var user entity.User

	if err := s.cache.Get(ctx, repository.BuildCacheKey(userCacheKey, uid), &user); err != nil {
		return entity.User{}, fmt.Errorf("userService - FindByID - s.cacheRepo.Get: %w", err)
	}

	if (user != entity.User{}) {
		return user, nil
	}

	user, err := s.userRepo.GetByID(ctx, uid)
	if err != nil {
		return entity.User{}, fmt.Errorf("userService - FindByID - s.userRepo.GetByID: %w", err)
	}

	if err = s.cache.Set(ctx, repository.BuildCacheKey(userCacheKey, uid), user, s.cacheTTL); err != nil {
		return entity.User{}, fmt.Errorf("userService - FindByID - s.cacheRepo.Set: %w", err)
	}

	return user, nil
}
func (s *userService) Register(ctx context.Context, email, password string) (entity.User, error) {
	bhash, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return entity.User{}, fmt.Errorf("userService - Register - bcrypt.GenerateFromPassword: %w", err)
	}

	dto, err := entity.NewUserSensitiveData(email, string(bhash))
	if err != nil {
		return entity.User{}, fmt.Errorf("userService - Register - entity.NewUserRegister: %w", err)
	}

	u, err := s.userRepo.Create(ctx, dto)
	if err != nil {
		return entity.User{}, fmt.Errorf("userService - Register - s.userRepo.Create: %w", err)
	}

	return u, nil
}

// Archive sets User isArchive state to archive
func (s *userService) Archive(ctx context.Context, uid string, archive bool) error {
	u, err := s.userRepo.Archive(ctx, uid, archive)
	if err != nil {
		return fmt.Errorf("userService - Archive - s.userRepo.Archive: %w", err)
	}

	if err = s.cache.Add(ctx, repository.BuildCacheKey(userCacheKey, uid), u, s.cacheTTL); err != nil {
		return fmt.Errorf("userService - Archive - s.cacheRepo.Add: %w", err)
	}

	return nil
}
