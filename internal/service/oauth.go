package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/domain"
	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/util"
)

type oauthService struct {
	cfg     *config.Config
	account Account
	session Session
}

func NewOAuthService(cfg *config.Config, a Account, s Session) *oauthService {
	return &oauthService{
		cfg:     cfg,
		account: a,
		session: s,
	}
}

func (s *oauthService) GetAuthorizeURI(ctx context.Context, provider string) (string, error) {
	provider = strings.ToLower(provider)

	scope, err := util.UniqueString(32)
	if err != nil {
		return "", fmt.Errorf("oauthService - GetAuthorizeURI - util.UniqueString: %w", err)
	}

	uri, err := url.Parse(s.cfg.Endpoints()[provider].AuthURL)
	if err != nil {
		return "", fmt.Errorf("oauthService - GetAuthorizeURI - url.Parse: %w", err)
	}

	q := uri.Query()
	q.Set("client_id", s.cfg.ClientIDs()[provider])
	q.Set("scope", s.cfg.Scopes()[provider])
	q.Set("state", scope)
	uri.RawQuery = q.Encode()

	return uri.String(), nil
}

func (s *oauthService) exchangeCode(ctx context.Context, provider, code string) (*oauth2.Token, error) {
	oauth := oauth2.Config{
		ClientID:     s.cfg.ClientIDs()[provider],
		ClientSecret: s.cfg.ClientSecrets()[provider],
		Endpoint:     s.cfg.Endpoints()[provider],
		Scopes:       strings.Split(s.cfg.Scopes()[provider], ","),
	}

	token, err := oauth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth.Exchange: %w", err)
	}

	return token, nil
}

func (s *oauthService) getGitHubUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	ts := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(ctx, ts)
	gh := github.NewClient(tc)

	u, r, err := gh.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("gh.Users.Get: %w", err)
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gh.Users.Get: %w", apperrors.ErrAuthGitHubUserNotReceived)
	}

	return u, nil
}

func (s *oauthService) createOrLogin(ctx context.Context, email string, d domain.Device) (domain.SessionCookie, error) {
	var aid string

	acc, err := s.account.GetByEmail(ctx, email)
	if err == nil {
		aid = acc.ID
	} else {
		if !errors.Is(err, apperrors.ErrAccountNotFound) {
			return domain.SessionCookie{}, fmt.Errorf("s.accountService.GetByEmail: %w", err)
		}

		aid, err = s.account.Create(ctx, email, util.RandomSpecialString(16))
		if err != nil {
			return domain.SessionCookie{}, fmt.Errorf("s.accountService.Create: %w", err)
		}
	}

	sess, err := s.session.Create(ctx, aid, d)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("s.sessionService.Create: %w", err)
	}

	return domain.NewSessionCookie(sess.ID, sess.TTL, &s.cfg.Session), nil
}

func (s *oauthService) GitHubLogin(ctx context.Context, code string, d domain.Device) (domain.SessionCookie, error) {
	token, err := s.exchangeCode(ctx, "github", code)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("oauthService - GitHubLogin - s.exchangeCode: %w", err)
	}

	u, err := s.getGitHubUser(ctx, token)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("oauthService - GitHubLogin - s.getGitHubUser: %w", err)
	}

	c, err := s.createOrLogin(ctx, *u.Email, d)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("oauthService - GitHubLogin - s.createOrLogin: %w", err)
	}

	return c, nil
}

func (s *oauthService) GoogleLogin(ctx context.Context, code string, d domain.Device) (domain.SessionCookie, error) {
	panic("implement")
	return domain.SessionCookie{}, nil
}
