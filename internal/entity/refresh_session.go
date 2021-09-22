package entity

import (
	"github.com/google/uuid"
	"time"
)

type RefreshSession struct {
	userAgent   string
	userIP      string
	fingerprint uuid.UUID

	RefreshToken uuid.UUID
	UserID       int
	ExpiresIn    time.Duration
	CreatedAt    time.Time
}

func NewRefreshSession(ua string, uip string, fp uuid.UUID) RefreshSession {
	return RefreshSession{
		userAgent:   ua,
		userIP:      uip,
		fingerprint: fp,
	}
}
