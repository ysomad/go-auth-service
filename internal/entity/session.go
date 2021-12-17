package entity

import (
	"encoding/json"
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
	UserID    string    `bson:"userID" redis:"userID"`
	UserAgent string    `bson:"userAgent" redis:"userAgent"`
	UserIP    string    `bson:"userIP" redis:"userIP"`
	TTL       int       `bson:"ttl" redis:"ttl"`
	ExpiresAt int64     `bson:"expiresAt" redis:"expiresAt"`
	CreatedAt time.Time `bson:"createdAt" redis:"createdAt"`
}

func NewSession(uid string, userAgent string, ip string, ttl time.Duration) (Session, error) {
	// TODO: add validation

	id, err := util.UniqueString(32)
	if err != nil {
		return Session{}, err
	}

	now := time.Now()
	return Session{
		ID:        id,
		UserID:    uid,
		UserAgent: userAgent,
		UserIP:    ip,
		TTL:       int(ttl.Seconds()),
		ExpiresAt: now.Add(ttl).Unix(),
		CreatedAt: now,
	}, nil
}

func (s Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
