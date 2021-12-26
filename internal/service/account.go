package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ysomad/go-auth-service/internal/domain"
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
	acc := domain.Account{Email: email}

	if err := acc.GeneratePasswordHash(password); err != nil {
		return fmt.Errorf("accountService - Create - acc.GeneratePasswordHash: %w", err)
	}

	if err := s.accountRepo.Create(ctx, acc); err != nil {
		return fmt.Errorf("accountService - Create - s.accountRepo.Create: %w", err)
	}

	return nil
}

func (s *accountService) GetByID(ctx context.Context, aid string) (domain.Account, error) {
	var acc domain.Account

	if err := s.cacheRepo.Get(ctx, aid, &acc); err == nil {
		return acc, nil
	}

	acc, err := s.accountRepo.FindByID(ctx, aid)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - GetByID - s.accountRepo.FindByID: %w", err)
	}

	// TODO: do not return on error when setting to cache
	if err = s.cacheRepo.Set(ctx, aid, acc, s.cacheTTL); err != nil {
		return domain.Account{}, fmt.Errorf("accountService - GetByID - s.cacheRepo.Set: %w", err)
	}

	return acc, nil
}

func (s *accountService) GetByEmail(ctx context.Context, email string) (domain.Account, error) {
	var acc domain.Account

	acc, err := s.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - GetByEmail - s.accountRepo.FindByEmail: %w", err)
	}

	return acc, nil
}

func (s *accountService) Archive(ctx context.Context, aid string) error {
	if err := s.cacheRepo.Delete(ctx, aid); err != nil {
		return fmt.Errorf("accountService - Archive - s.cacheRepo.Delete: %w", err)
	}

	if err := s.accountRepo.Archive(ctx, aid, true); err != nil {
		return fmt.Errorf("accountService - Archive - s.accountRepo.Archive: %w", err)
	}

	return nil
}
