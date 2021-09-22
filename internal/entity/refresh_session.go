package entity

import (
	"github.com/google/uuid"
	"time"
)

type RefreshSession struct {
	RefreshToken uuid.UUID
	UserID       int
	UserAgent    string
	UserIP       string
	Fingerprint  uuid.UUID
	ExpiresIn    time.Duration
	CreatedAt    time.Time
}

func (s *RefreshSession) SetUserID(id int) {
	s.UserID = id
}

func (s *RefreshSession) SetRefreshToken(t uuid.UUID) {
	s.RefreshToken = t
}

func (s *RefreshSession) SetExpiresIn(e time.Duration) {
	s.ExpiresIn = e
}

