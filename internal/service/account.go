package service

import (
	"context"
	"fmt"

	"github.com/ysomad/go-auth-service/internal/domain"
)

type accountService struct {
	accountRepo AccountRepo
}

func NewAccountService(r AccountRepo) *accountService {
	return &accountService{
		accountRepo: r,
	}
}

func (s *accountService) Create(ctx context.Context, email, password string) (string, error) {
	acc := domain.Account{Email: email}

	if err := acc.GeneratePasswordHash(password); err != nil {
		return "", fmt.Errorf("accountService - Create - acc.GeneratePasswordHash: %w", err)
	}

	aid, err := s.accountRepo.Create(ctx, acc)
	if err != nil {
		return "", fmt.Errorf("accountService - Create - s.accountRepo.Create: %w", err)
	}

	return aid, nil
}

func (s *accountService) GetByID(ctx context.Context, aid string) (domain.Account, error) {
	var acc domain.Account

	acc, err := s.accountRepo.FindByID(ctx, aid)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - GetByID - s.accountRepo.FindByID: %w", err)
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
	if err := s.accountRepo.Archive(ctx, aid, true); err != nil {
		return fmt.Errorf("accountService - Archive - s.accountRepo.Archive: %w", err)
	}

	return nil
}
