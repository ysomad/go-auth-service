package entity

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrSessionExpired = errors.New("session expired")
)

// Session represents refresh token session for JWT authentication
type Session struct {
	RefreshToken uuid.UUID     `json:"refreshToken" redis:"refreshToken"`
	UserID       uuid.UUID     `json:"userID" redis:"userID"`
	UserAgent    string        `json:"userAgent" redis:"userAgent"`
	UserIP       string        `json:"userIP" redis:"userIP"`
	Fingerprint  uuid.UUID     `json:"fingerprint" redis:"fingerprint"`
	ExpiresIn    time.Duration `json:"expiresIn" redis:"expiresIn"`
	ExpiresAt    int64         `json:"expiresAt" redis:"expiresAt"`
	CreatedAt    time.Time     `json:"createdAt" redis:"createdAt"`
}

func (s Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// SessionSecurityDTO stores session security data
type SessionSecurityDTO struct {
	RefreshToken uuid.UUID
	Fingerprint  uuid.UUID
	UserAgent    string
	UserIP       string
}
