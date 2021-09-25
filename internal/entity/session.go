package entity

import (
	"github.com/google/uuid"
	"time"
)

// Session represents refresh token session for JWT authentication
type Session struct {
	RefreshToken uuid.UUID     `json:"refreshToken" redis:"token"`
	UserID       uuid.UUID     `json:"userID" redis:"uid"`
	UserAgent    string        `json:"userAgent" redis:"ua"`
	UserIP       string        `json:"userIP" redis:"ip"`
	Fingerprint  uuid.UUID     `json:"fingerprint" redis:"fp"`
	ExpiresIn    time.Duration `json:"expiresIn"`
	ExpiresAt    int64         `json:"expiresAt" redis:"exp"`
	CreatedAt    time.Time     `json:"createdAt" redis:"created"`
}

type SessionSecurityDTO struct {
	RefreshToken uuid.UUID
	Fingerprint  uuid.UUID
	UserAgent    string
	UserIP       string
}