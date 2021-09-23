package entity

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// Session represents refresh token session for JWT authentication
type Session struct {
	RefreshToken uuid.UUID     `json:"refreshToken"`
	UserID       int           `json:"userId"`
	UserAgent    string        `json:"userAgent"`
	UserIP       string        `json:"userIP"`
	Fingerprint  uuid.UUID     `json:"fingerprint"`
	ExpiresAt    time.Time     `json:"expiresAt"`
	ExpiresIn    time.Duration `json:"expiresIn"`
	CreatedAt    time.Time     `json:"createdAt"`
}

func (s *Session) SetExpiresAt() error {
	if s.ExpiresIn == 0 || s.CreatedAt.IsZero() {
		return errors.New("session creation error")
	}

	s.ExpiresAt = s.CreatedAt.Add(s.ExpiresIn)

	return nil
}
