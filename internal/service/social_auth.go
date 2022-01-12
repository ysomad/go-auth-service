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

type socialAuthService struct {
	cfg            *config.Config
	accountService Account
	sessionService Session
}

func NewSocialAuthService(cfg *config.Config, a Account, s Session) *socialAuthService {
	return &socialAuthService{
		cfg:            cfg,
		accountService: a,
		sessionService: s,
	}
}

func (s *socialAuthService) GetAuthorizeURI(ctx context.Context, provider string) (string, error) {
	provider = strings.ToLower(provider)

	scope, err := util.UniqueString(32)
	if err != nil {
		return "", fmt.Errorf("socialAuthService - GetAuthorizeURI - util.UniqueString: %w", err)
	}

	uri, err := url.Parse(s.cfg.Endpoints()[provider].AuthURL)
	if err != nil {
		return "", fmt.Errorf("socialAuthService - GetAuthorizeURI - url.Parse: %w", err)
	}

	q := uri.Query()
	q.Set("client_id", s.cfg.ClientIDs()[provider])
	q.Set("scope", s.cfg.Scopes()[provider])
	q.Set("state", scope)
	uri.RawQuery = q.Encode()

	return uri.String(), nil
}

func (s *socialAuthService) GitHubLogin(ctx context.Context, code string, d domain.Device) (domain.SessionCookie, error) {
	token, err := s.exchangeCode(ctx, domain.ProviderGitHub, code)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("socialAuthService  - GitHubLogin - s.exchangeCode: %w", err)
	}

	u, err := s.getGitHubUser(ctx, token)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("socialAuthService - GitHubLogin - s.getGitHubUser: %w", err)
	}

	c, err := s.loginOrSignUp(ctx, *u.Email, *u.Login, domain.ProviderGitHub, d)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("socialAuthService - GitHubLogin - s.createOrLogin: %w", err)
	}

	return c, nil
}

func (s *socialAuthService) GoogleLogin(ctx context.Context, code string, d domain.Device) (domain.SessionCookie, error) {
	panic("implement")
	return domain.SessionCookie{}, nil
}

// --------------------------------------------------------------------------------------------------------------------

// exchangeCode sends OAuth2 authorization code to data provider authorization server in order to
// get REST API access token which is used to use private provider api.
func (s *socialAuthService) exchangeCode(ctx context.Context, provider, code string) (*oauth2.Token, error) {
	o := oauth2.Config{
		ClientID:     s.cfg.ClientIDs()[provider],
		ClientSecret: s.cfg.ClientSecrets()[provider],
		Endpoint:     s.cfg.Endpoints()[provider],
		Scopes:       strings.Split(s.cfg.Scopes()[provider], ","),
	}

	t, err := o.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("o.Exchange: %w", err)
	}

	return t, nil
}

// getGitHubUser returns github user using access token received from exchangeCode method.
func (s *socialAuthService) getGitHubUser(ctx context.Context, t *oauth2.Token) (*github.User, error) {
	ts := oauth2.StaticTokenSource(t)
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

// loginOrSignUp logs in user with received data from OAuth2 data provider if account exist or creates
// new account with random password and logs it in.
func (s *socialAuthService) loginOrSignUp(ctx context.Context, email, username,
	provider string, d domain.Device) (domain.SessionCookie, error) {

	var aid string

	acc, err := s.accountService.GetByEmail(ctx, email)
	if err == nil {
		aid = acc.ID
	} else {
		if !errors.Is(err, apperrors.ErrAccountNotFound) {
			return domain.SessionCookie{}, fmt.Errorf("s.accountService.GetByEmail: %w", err)
		}

		acc = domain.Account{Email: email, Username: username}
		acc.RandomPassword()

		aid, err = s.accountService.Create(ctx, acc)
		if err != nil {
			return domain.SessionCookie{}, fmt.Errorf("s.accountService.Create: %w", err)
		}
	}

	sess, err := s.sessionService.Create(ctx, aid, provider, d)
	if err != nil {
		return domain.SessionCookie{}, fmt.Errorf("s.sessionService.Create: %w", err)
	}

	return domain.NewSessionCookie(sess.ID, sess.TTL, &s.cfg.Session), nil
}
