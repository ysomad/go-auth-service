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
	AccountID string    `json:"accountId" bson:"accountId"`
	Provider  string    `json:"provider" bson:"provider"`
	UserAgent string    `json:"userAgent" bson:"userAgent"`
	IP        string    `json:"ip" bson:"ip"`
	TTL       int       `json:"ttl" bson:"ttl"`
	ExpiresAt int64     `json:"expiresAt" bson:"expiresAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

func NewSession(aid, provider, userAgent, ip string, ttl time.Duration) (Session, error) {
	id, err := util.UniqueString(32)
	if err != nil {
		return Session{}, fmt.Errorf("utils.UniqueString: %w", apperrors.ErrSessionNotCreated)
	}

	now := time.Now()

	return Session{
		ID:        id,
		AccountID: aid,
		Provider:  provider,
		UserAgent: userAgent,
		IP:        ip,
		TTL:       int(ttl.Seconds()),
		ExpiresAt: now.Add(ttl).Unix(),
		CreatedAt: now,
	}, nil
}
