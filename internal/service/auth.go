package service

import (
	"context"
	"fmt"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/auth"
)

type authService struct {
	accountService Account
	sessionService Session
	tokenManager   auth.TokenManager
}

func NewAuthService(a Account, s Session, t auth.TokenManager) *authService {
	return &authService{
		accountService: a,
		sessionService: s,
		tokenManager:   t,
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

func (s *authService) NewAccessToken(ctx context.Context, aid, password string) (domain.Token, error) {
	acc, err := s.accountService.GetByID(ctx, aid)
	if err != nil {
		return domain.Token{}, fmt.Errorf("authService - NewAccessToken - s.accountService.GetByID: %w", err)
	}

	if err := acc.CompareHashAndPassword(password); err != nil {
		return domain.Token{}, fmt.Errorf("authService - NewAccessToken - acc.CompareHashAndPassword: %w", err)
	}

	token, err := s.tokenManager.NewJWT(aid)
	if err != nil {
		return domain.Token{}, fmt.Errorf("authService - NewAccessToken - s.tokenManager.NewJWT: %w", err)
	}

	return domain.Token{
		AccessToken: token,
	}, nil
}

func (s *authService) ParseAccessToken(ctx context.Context, token string) (string, error) {
	aid, err := s.tokenManager.ParseJWT(token)
	if err != nil {
		return "", fmt.Errorf("authService - ParseAccessToken - s.tokenManager.ParseJWT: %w", err)
	}

	return aid, nil
}
