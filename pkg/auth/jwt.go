package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt"
)

type JWT interface {
	NewAccess(userID int) (string, error)
	NewRefresh() (uuid.UUID, error)
}

type JWTManager struct {
	signingKey     string
	accessTokenTTL time.Duration
}

func NewJWTManager(signingKey string, accessTokenTTL time.Duration) (*JWTManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &JWTManager{
		signingKey:     signingKey,
		accessTokenTTL: accessTokenTTL,
	}, nil
}

func (m *JWTManager) NewAccess(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(m.accessTokenTTL).Unix(),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *JWTManager) NewRefresh() (uuid.UUID, error) {
	return uuid.NewRandom()
}
