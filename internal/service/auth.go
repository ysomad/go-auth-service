package service

import (
	"context"
	"fmt"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/jwt"
)

type authService struct {
	cfg     *config.Config
	token   jwt.Token
	account Account
	session Session
}

func NewAuthService(cfg *config.Config, t jwt.Token, a Account, s Session) *authService {
	return &authService{
		cfg:     cfg,
		token:   t,
		account: a,
		session: s,
	}
}

func (s *authService) EmailLogin(ctx context.Context, email, password string, d Device) (domain.Session, error) {
	a, err := s.account.GetByEmail(ctx, email)
	if err != nil {
		return domain.Session{}, fmt.Errorf("authService - EmailLogin - s.account.GetByEmail: %w", err)
	}

	a.Password = password

	if err = a.CompareHashAndPassword(); err != nil {
		return domain.Session{}, fmt.Errorf("authService - EmailLogin - a.CompareHashAndPassword: %w", err)
	}

	sess, err := s.session.Create(ctx, a.ID, providerEmail, d)
	if err != nil {
		return domain.Session{}, fmt.Errorf("authService - EmailLogin - s.session.Create: %w", err)
	}

	return sess, nil

}

func (s *authService) Logout(ctx context.Context, sid string) error {
	if err := s.session.Terminate(ctx, sid, ""); err != nil {
		return fmt.Errorf("authService - Logout - s.session.Terminate: %w", err)
	}

	return nil
}

func (s *authService) NewAccessToken(ctx context.Context, aid, password string) (string, error) {
	a, err := s.account.GetByID(ctx, aid)
	if err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - s.account.GetByID: %w", err)
	}

	a.Password = password

	if err := a.CompareHashAndPassword(); err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - a.CompareHashAndPassword: %w", err)
	}

	t, err := s.token.New(aid)
	if err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - s.token.New: %w", err)
	}

	return t, nil
}

func (s *authService) ParseAccessToken(ctx context.Context, t string) (string, error) {
	aid, err := s.token.Parse(t)
	if err != nil {
		return "", fmt.Errorf("authService - ParseAccessToken - s.token.Parse: %w", err)
	}

	return aid, nil
}
