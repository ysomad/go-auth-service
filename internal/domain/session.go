package domain

import (
	"errors"
	"time"

	"github.com/ysomad/go-auth-service/pkg/util"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrUnauthorized    = errors.New("unauthorized")
)

// Session represents refresh token session for JWT authentication
type Session struct {
	ID        string    `bson:"_id" redis:"id"`
	AccountID string    `bson:"accountID" redis:"accountID"`
	UserAgent string    `bson:"userAgent" redis:"userAgent"`
	ClientIP  string    `bson:"clientIP" redis:"clientIP"`
	TTL       int       `bson:"ttl" redis:"ttl"`
	ExpiresAt int64     `bson:"expiresAt" redis:"expiresAt"`
	CreatedAt time.Time `bson:"createdAt" redis:"createdAt"`
}

func NewSession(aid string, userAgent string, ip string, ttl time.Duration) (Session, error) {
	// TODO: add validation

	id, err := util.UniqueString(32)
	if err != nil {
		// TODO: return generic err pkg/httperror
		return Session{}, err
	}

	now := time.Now()
	return Session{
		ID:        id,
		AccountID: aid,
		UserAgent: userAgent,
		ClientIP:  ip,
		TTL:       int(ttl.Seconds()),
		ExpiresAt: now.Add(ttl).Unix(),
		CreatedAt: now,
	}, nil
}

/*
func (s Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
*/

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
