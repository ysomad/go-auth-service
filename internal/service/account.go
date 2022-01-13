package service

import (
	"context"
	"fmt"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/domain"
)

type accountService struct {
	cfg *config.Config
	repo AccountRepo
	session Session
}

func NewAccountService(cfg *config.Config, r AccountRepo, s Session) *accountService {
	return &accountService{
		cfg: cfg,
		repo: r,
		session: s,
	}
}

func (s *accountService) Create(ctx context.Context, a domain.Account) (string, error) {
	if err := a.GeneratePasswordHash(); err != nil {
		return "", fmt.Errorf("accountService - Create - acc.GeneratePasswordHash: %w", err)
	}

	aid, err := s.repo.Create(ctx, a)
	if err != nil {
		return "", fmt.Errorf("accountService - Create - s.repo.Create: %w", err)
	}

	return aid, nil
}

func (s *accountService) GetByID(ctx context.Context, aid string) (domain.Account, error) {
	var acc domain.Account

	acc, err := s.repo.FindByID(ctx, aid)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - GetByID - s.repo.FindByID: %w", err)
	}

	return acc, nil
}

func (s *accountService) GetByEmail(ctx context.Context, email string) (domain.Account, error) {
	var acc domain.Account

	acc, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - GetByEmail - s.repo.FindByEmail: %w", err)
	}

	return acc, nil
}

func (s *accountService) Delete(ctx context.Context, aid, sid string) (SessionCookie, error) {
	if err := s.repo.Archive(ctx, aid, true); err != nil {
		return SessionCookie{}, fmt.Errorf("accountService - Archive - s.repo.Archive: %w", err)
	}

	if err := s.session.TerminateAll(ctx, aid, sid); err != nil {
		return SessionCookie{}, fmt.Errorf("accountService - Archive - s.session.TerminateAll: %w", err)
	}

	return NewSessionCookie(sid, -1, &s.cfg.Session), nil
}

func (s *accountService) Verify(ctx context.Context, code string) error {
	panic("implement")

	return nil
}
