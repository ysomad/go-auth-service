package domain

import (
	"fmt"
	"time"

	"github.com/ysomad/go-auth-service/config"
	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/util"
)

// Provider constants to track how user is logged in
const (
	ProviderEmail    = "email"
	ProviderUsername = "username"
	ProviderGitHub   = "github"
	ProviderGoogle   = "google"
)

// Session represents refresh token session for JWT authentication
type Session struct {
	ID        string    `json:"id" bson:"_id"`
	AccountID string    `json:"accountID" bson:"accountID"`
	Provider  string    `json:"provider" bson:"provider"` // TODO: store provider in session
	UserAgent string    `json:"userAgent" bson:"userAgent"`
	IP        string    `json:"ip" bson:"ip"`
	TTL       int       `json:"ttl" bson:"ttl"`
	ExpiresAt int64     `json:"expiresAt" bson:"expiresAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

func NewSession(aid, provider, userAgent, ip string, ttl time.Duration) (Session, error) {
	id, err := util.UniqueString(32)
	if err != nil {
		return Session{}, fmt.Errorf("utils.UniqueString: %w", apperrors.ErrSessionNotCreated)
	}

	now := time.Now()

	return Session{
		ID:        id,
		AccountID: aid,
		Provider:  provider,
		UserAgent: userAgent,
		IP:        ip,
		TTL:       int(ttl.Seconds()),
		ExpiresAt: now.Add(ttl).Unix(),
		CreatedAt: now,
	}, nil
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
