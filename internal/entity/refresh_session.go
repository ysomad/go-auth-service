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
	ExpiresAt    int64
	ExpiresIn    time.Duration
	CreatedAt    time.Time
}
