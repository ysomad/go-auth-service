package domain

import (
	"fmt"
	"time"

	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/util"
)

// Session represents refresh token session for JWT authentication
type Session struct {
	ID        string    `json:"id" bson:"_id"`
	AccountID string    `json:"accountID" bson:"accountID"`
	UserAgent string    `json:"userAgent" bson:"userAgent"`
	IP        string    `json:"ip" bson:"ip"`
	TTL       int       `json:"ttl" bson:"ttl"`
	ExpiresAt int64     `json:"expiresAt" bson:"expiresAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

func NewSession(aid string, userAgent string, ip string, ttl time.Duration) (Session, error) {
	id, err := util.UniqueString(32)
	if err != nil {
		return Session{}, fmt.Errorf("utils.UniqueString: %w", apperrors.ErrSessionNotCreated)
	}

	now := time.Now()

	return Session{
		ID:        id,
		AccountID: aid,
		UserAgent: userAgent,
		IP:        ip,
		TTL:       int(ttl.Seconds()),
		CreatedAt: now,
	}, nil
}

// SessionCookie represents data transfer object which
// contains data needed to create a cookie.
type SessionCookie struct {
	id  string
	ttl int
}

func NewSessionCookie(sid string, ttl int) SessionCookie {
	return SessionCookie{
		id:  sid,
		ttl: ttl,
	}
}

func (s SessionCookie) ID() string { return s.id }
func (s SessionCookie) TTL() int   { return s.ttl }
