package service

import (
	"context"
	"fmt"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/domain"
)

type sessionService struct {
	cfg  *config.Config
	repo SessionRepo
}

func NewSessionService(cfg *config.Config, s SessionRepo) *sessionService {
	return &sessionService{
		cfg:  cfg,
		repo: s,
	}
}

func (s *sessionService) Create(ctx context.Context, aid, provider string, d domain.Device) (domain.Session, error) {
	sess, err := domain.NewSession(aid, provider, d.UserAgent, d.IP, s.cfg.Session.TTL)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Create - domain.NewSession: %w", err)
	}

	if err = s.repo.Create(ctx, sess); err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Create - s.sessionRepo.Create: %w", err)
	}

	return sess, nil
}

func (s *sessionService) Get(ctx context.Context, sid string) (domain.Session, error) {
	var sess domain.Session

	sess, err := s.repo.Get(ctx, sid)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Get - s.sessionRepo.Get: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetAll(ctx context.Context, aid string) ([]domain.Session, error) {
	sessions, err := s.repo.GetAll(ctx, aid)
	if err != nil {
		return nil, fmt.Errorf("sessionService - GetAll - s.sessionRepo.GetAll: %w", err)
	}

	return sessions, nil
}

func (s *sessionService) Terminate(ctx context.Context, sid string) error {
	if err := s.repo.Delete(ctx, sid); err != nil {
		return fmt.Errorf("sessionService - Terminate - s.sessionRepo.Delete: %w", err)
	}

	return nil
}

func (s *sessionService) TerminateAll(ctx context.Context, aid, currSid string) error {
	if err := s.repo.DeleteAll(ctx, aid, currSid); err != nil {
		return fmt.Errorf("sessionService - TerminateAll - s.sessionRepo.DeleteAll: %w", err)
	}

	return nil
}
