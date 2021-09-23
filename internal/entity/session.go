package entity

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// Session represents refresh token session for JWT authentication
type Session struct {
	RefreshToken uuid.UUID
	UserID       int
	UserAgent    string
	UserIP       string
	Fingerprint  uuid.UUID
	ExpiresAt    time.Time
	ExpiresIn    time.Duration
	CreatedAt    time.Time
}

func (s *Session) SetExpiresAt() error {
	if s.ExpiresIn == 0 || s.CreatedAt.IsZero()  {
		return errors.New("session creation error")
	}

	s.ExpiresAt = s.CreatedAt.Add(s.ExpiresIn)

	return nil
}