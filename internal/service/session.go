package service

import (
	"context"
	"fmt"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/apperrors"
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

// SessionCookie represents data transfer object which
// contains data needed to create a cookie.
type SessionCookie struct {
	ID       string
	TTL      int
	Domain   string
	Secure   bool
	HTTPOnly bool
	Key      string
}

func NewSessionCookie(sid string, ttl int, cfg *config.Session) SessionCookie {
	return SessionCookie{
		ID:       sid,
		TTL:      ttl,
		Domain:   cfg.CookieDomain,
		Secure:   cfg.CookieSecure,
		HTTPOnly: cfg.CookieHTTPOnly,
		Key:      cfg.CookieKey,
	}
}

// Device represents data transfer object with user device data
type Device struct {
	UserAgent string
	IP        string
}

func NewDevice(ua string, ip string) Device {
	return Device{
		UserAgent: ua,
		IP:        ip,
	}
}

func (s *sessionService) Create(ctx context.Context, aid, provider string, d Device) (domain.Session, error) {
	sess, err := domain.NewSession(aid, provider, d.UserAgent, d.IP, s.cfg.Session.TTL)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Create - domain.NewSession: %w", err)
	}

	if err = s.repo.Create(ctx, sess); err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Create - s.repo.Create: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetByID(ctx context.Context, sid string) (domain.Session, error) {
	sess, err := s.repo.FindByID(ctx, sid)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Get - s.repo.FindByID: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetAll(ctx context.Context, aid string) ([]domain.Session, error) {
	sessions, err := s.repo.FindAll(ctx, aid)
	if err != nil {
		return nil, fmt.Errorf("sessionService - GetAll - s.repo.FindAll: %w", err)
	}

	return sessions, nil
}

func (s *sessionService) Terminate(ctx context.Context, sid, currSid string) error {
	if sid == currSid {
		return fmt.Errorf("sessionService - Terminate: %w", apperrors.ErrSessionNotTerminated)
	}

	if err := s.repo.Delete(ctx, sid); err != nil {
		return fmt.Errorf("sessionService - Terminate - s.sessionRepo.Delete: %w", err)
	}

	return nil
}

func (s *sessionService) TerminateAll(ctx context.Context, aid, sid string) error {
	if err := s.repo.DeleteAll(ctx, aid, sid); err != nil {
		return fmt.Errorf("sessionService - TerminateAll - s.repo.DeleteAll: %w", err)
	}

	return nil
}
