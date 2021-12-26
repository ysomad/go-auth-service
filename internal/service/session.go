package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ysomad/go-auth-service/internal/domain"
)

type sessionService struct {
	accountRepo AccountRepo

	sessionRepo SessionRepo
	sessionTTL  time.Duration
}

func NewSessionService(a AccountRepo, s SessionRepo, sessionTTL time.Duration) *sessionService {

	return &sessionService{
		accountRepo: a,
		sessionRepo: s,
		sessionTTL:  sessionTTL,
	}
}

func (s *sessionService) Create(ctx context.Context, aid string, d domain.Device) (domain.Session, error) {
	sess, err := domain.NewSession(aid, d.UserAgent, d.IP, s.sessionTTL)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Create - domain.NewSession: %w", err)
	}

	if err = s.sessionRepo.Create(ctx, sess); err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Create - s.sessionRepo.Create: %w", err)
	}

	return sess, nil
}

func (s *sessionService) Get(ctx context.Context, sid string) (domain.Session, error) {
	var sess domain.Session

	sess, err := s.sessionRepo.Get(ctx, sid)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Get - s.sessionRepo.Get: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetAll(ctx context.Context, aid string) ([]domain.Session, error) {
	sessions, err := s.sessionRepo.GetAll(ctx, aid)
	if err != nil {
		return nil, fmt.Errorf("sessionService - GetAll - s.sessionRepo.GetAll: %w", err)
	}

	return sessions, nil
}

func (s *sessionService) Terminate(ctx context.Context, sid string) error {
	if err := s.sessionRepo.Delete(ctx, sid); err != nil {
		return fmt.Errorf("sessionService - Terminate - s.sessionRepo.Delete: %w", err)
	}

	return nil
}

func (s *sessionService) TerminateAll(ctx context.Context, uid string) error {
	if err := s.sessionRepo.DeleteAll(ctx, uid); err != nil {
		return fmt.Errorf("sessionService - TerminateAll - s.sessionRepo.DeleteAll: %w", err)
	}

	return nil
}
