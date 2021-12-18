package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service/repository"
)

type accountService struct {
	accountRepo AccountRepo

	cacheRepo CacheRepo
	cacheTTL  time.Duration
}

func NewAccountService(r AccountRepo, c CacheRepo, cacheTTL time.Duration) *accountService {
	return &accountService{
		accountRepo: r,
		cacheRepo:   c,
		cacheTTL:    cacheTTL,
	}
}

func (s *accountService) Create(ctx context.Context, email, password string) error {
	// TODO: refactor
	// TODO: return generic err from httperror
	bhash, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return fmt.Errorf("userService - Register - bcrypt.GenerateFromPassword: %w", err)
	}

	cred, err := domain.NewAccountCredentials(email, string(bhash))
	if err != nil {
		return fmt.Errorf("userService - Register - entity.NewUserRegister: %w", err)
	}

	if err := s.accountRepo.Create(ctx, cred); err != nil {
		return fmt.Errorf("userService - Register - s.userRepo.Create: %w", err)
	}

	return nil
}

func (s *accountService) GetByID(ctx context.Context, aid string) (domain.Account, error) {
	// TODO: RETURN GENERIC ERR FROM pkg httperror
	// TODO: refactor
	var acc domain.Account

	if err := s.cacheRepo.Get(ctx, repository.BuildCacheKey("acc", aid), &acc); err != nil {
		return domain.Account{}, fmt.Errorf("userService - FindByID - s.cacheRepo.Get: %w", err)
	}

	if (acc != domain.Account{}) {
		return acc, nil
	}

	acc, err := s.accountRepo.FindByID(ctx, aid)
	if err != nil {
		return domain.Account{}, fmt.Errorf("userService - FindByID - s.userRepo.GetByID: %w", err)
	}

	if err = s.cacheRepo.Set(ctx, repository.BuildCacheKey("acc", aid), acc, s.cacheTTL); err != nil {
		return domain.Account{}, fmt.Errorf("userService - FindByID - s.cacheRepo.Set: %w", err)
	}

	return acc, nil
}

func (s *accountService) GetByEmail(ctx context.Context, email string) (domain.Account, error) {
	// TODO: RETURN GENERIC ERR FROM pkg httperror
	// TODO: refactor
	var acc domain.Account

	if err := s.cacheRepo.Get(ctx, repository.BuildCacheKey("acc", email), &acc); err != nil {
		return domain.Account{}, fmt.Errorf("userService - FindByID - s.cacheRepo.Get: %w", err)
	}

	if (acc != domain.Account{}) {
		return acc, nil
	}

	acc, err := s.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return domain.Account{}, fmt.Errorf("userService - FindByID - s.userRepo.GetByID: %w", err)
	}

	if err = s.cacheRepo.Set(ctx, repository.BuildCacheKey("acc", email), acc, s.cacheTTL); err != nil {
		return domain.Account{}, fmt.Errorf("userService - FindByID - s.cacheRepo.Set: %w", err)
	}

	return acc, nil
}

func (s *accountService) Archive(ctx context.Context, aid string) error {
	// TODO: refactor
	// TODO: return generic err from httperror

	if err := s.accountRepo.Archive(ctx, aid, true); err != nil {
		return fmt.Errorf("userService - Archive - s.userRepo.Archive: %w", err)
	}

	if err := s.cacheRepo.Delete(ctx, repository.BuildCacheKey("acc", aid)); err != nil {
		return fmt.Errorf("userService - Archive - s.cacheRepo.Add: %w", err)
	}

	return nil
}
