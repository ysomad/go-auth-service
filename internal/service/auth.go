package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/pkg/auth"
	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/util"
)

type authService struct {
	cfg            config.Auth
	accountService Account
	sessionService Session
	jwtManager     auth.JWTManager
}

func NewAuthService(cfg config.Auth, a Account, s Session, m auth.JWTManager) *authService {
	return &authService{
		cfg:            cfg,
		accountService: a,
		sessionService: s,
		jwtManager:     m,
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

func (s *authService) NewAccessToken(ctx context.Context, aid, password string) (string, error) {
	acc, err := s.accountService.GetByID(ctx, aid)
	if err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - s.accountService.GetByID: %w", err)
	}

	if err := acc.CompareHashAndPassword(password); err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - acc.CompareHashAndPassword: %w", err)
	}

	token, err := s.jwtManager.New(aid)
	if err != nil {
		return "", fmt.Errorf("authService - NewAccessToken - s.tokenManager.NewJWT: %w", err)
	}

	return token, nil
}

func (s *authService) ParseAccessToken(ctx context.Context, token string) (string, error) {
	aid, err := s.jwtManager.Parse(token)
	if err != nil {
		return "", fmt.Errorf("authService - ParseAccessToken - s.tokenManager.ParseJWT: %w", err)
	}

	return aid, nil
}

func (s *authService) GitHubLogin(ctx context.Context, code string, d domain.Device) (domain.SessionCookie, error) {
	// Request access token from github using code
	req, err := http.NewRequest("POST", oauth2github.Endpoint.TokenURL, nil)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - GitHubLogin - http.NewRequest: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Set("client_id", s.cfg.GitHubClientID)
	q.Set("client_secret", s.cfg.GitHubClientSecret)
	q.Set("code", code)
	req.URL.RawQuery = q.Encode()

	c := new(http.Client)

	resp, err := c.Do(req)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - GitHubLogin - c.Do: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - GitHubLogin - ioutil.ReadAll: %w", err)
	}

	var token oauth2.Token

	if err := json.Unmarshal(body, &token); err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - GitHubLogin - json.Unmarshal: %w", err)
	}

	ts := oauth2.StaticTokenSource(&token)
	tc := oauth2.NewClient(ctx, ts)

	gh := github.NewClient(tc)

	ghu, ghr, err := gh.Users.Get(ctx, "")
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - GitHubLogin - gh.Users.Get: %w", err)
	}

	if ghr.StatusCode != http.StatusOK {
		return domain.SessionCookie{}, fmt.Errorf("authService - GitHubLogin - gh.Users.Get: %w", apperrors.ErrAuthGitHubUserNotReceived)
	}

	// refactor ???
	var aid string

	acc, err := s.accountService.GetByEmail(ctx, *ghu.Email)
	if err == nil {
		aid = acc.ID
	} else {
		if !errors.Is(err, apperrors.ErrAccountNotFound) {
			return domain.SessionCookie{}, fmt.Errorf("authService - GitHubLogin - s.accountService.GetByEmail: %w", err)
		}

		aid, err = s.accountService.Create(ctx, *ghu.Email, util.RandomSpecialString(16))
		if err != nil {
			return domain.SessionCookie{}, fmt.Errorf("authService - GitHubLogin - s.accountService.Create: %w", err)
		}
	}

	sess, err := s.sessionService.Create(ctx, aid, d)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("authService - GitHubLogin - s.sessionService.Create: %w", err)
	}

	return domain.NewSessionCookie(sess.ID, sess.TTL), nil
}
