package service

import (
	"context"
	"fmt"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/auth"
)

type authService struct {
	cfg     *config.Config
	jwt     auth.JWTManager
	account Account
	session Session
}

func NewAuthService(cfg *config.Config, jwt auth.JWTManager, a Account, s Session) *authService {
	return &authService{
		cfg:     cfg,
		jwt:     jwt,
		account: a,
		session: s,
	}
}

func (s *authService) EmailLogin(ctx context.Context, email, password string, d domain.Device) (domain.SessionCookie, error) {
	acc, err := s.account.GetByEmail(ctx, email)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - EmailLogin - s.accountService.GetByEmail: %w", err)
	}

	acc.Password = password

	if err = acc.CompareHashAndPassword(); err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - EmailLogin - acc.CompareHashAndPassword: %w", err)
	}

	sess, err := s.session.Create(ctx, acc.ID, domain.ProviderEmail, d)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - EmailLogin - s.sessionService.Create: %w", err)
	}

	return domain.NewSessionCookie(sess.ID, sess.TTL, &s.cfg.Session), nil
}

func (s *authService) Logout(ctx context.Context, sid string) error {
	if err := s.session.Terminate(ctx, sid); err != nil {
		return fmt.Errorf("authService - Logout - s.sessionService.Terminate: %w", err)
	}

	return nil
}

func (s *authService) NewAccessToken(ctx context.Context, aid, password string) (string, error) {
	acc, err := s.account.GetByID(ctx, aid)
	if err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - s.accountService.GetByID: %w", err)
	}

	acc.Password = password

	if err := acc.CompareHashAndPassword(); err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - acc.CompareHashAndPassword: %w", err)
	}

	token, err := s.jwt.New(aid)
	if err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - s.tokenManager.NewJWT: %w", err)
	}

	return token, nil
}

func (s *authService) ParseAccessToken(ctx context.Context, token string) (string, error) {
	aid, err := s.jwt.Parse(token)
	if err != nil {
		return "", fmt.Errorf("authService - ParseAccessToken - s.tokenManager.ParseJWT: %w", err)
	}

	return aid, nil
}
