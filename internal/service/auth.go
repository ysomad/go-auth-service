package service

import (
	"context"
	"fmt"

	"github.com/ysomad/go-auth-service/internal/domain"
)

type authService struct {
	accountService Account
	sessionService Session
}

func NewAuthService(a Account, s Session) *authService {
	return &authService{
		accountService: a,
		sessionService: s,
	}
}

func (s *authService) EmailLogin(ctx context.Context, email, password string, d domain.Device) (domain.SessionCookie, error) {
	acc, err := s.accountService.GetByEmail(ctx, email)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - EmailLogin - s.accountService.GetByEmail: %w", err)
	}

	if err = acc.CompareHashAndPassword(password); err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - EmailLogin - acc.CompareHashAndPassword: %w", err)
	}

	sess, err := s.sessionService.Create(ctx, acc.ID, d)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - EmailLogin - s.sessionService.Create: %w", err)
	}

	return domain.NewSessionCookie(sess.ID, sess.TTL), nil
}

func (s *authService) Logout(ctx context.Context, sid string) error {
	if err := s.sessionService.Terminate(ctx, sid); err != nil {
		return fmt.Errorf("authService - Logout - s.sessionService.Terminate: %w", err)
	}

	return nil
}

func (s *authService) GetAccessToken(ctx context.Context, aid string) (domain.Token, error) {
  panic("implement")
	return domain.Token{}, nil
}
