package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrNoSigningKey         = errors.New("empty signing key")
	ErrNoClaims             = errors.New("error getting claims from token")
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
)

type TokenManager interface {
	NewJWT(sub string) (string, error)
	ParseJWT(accessToken string) (string, error)
}

type tokenManager struct {
	signingKey string
	ttl        time.Duration
}

func NewTokenManager(signingKey string, ttl time.Duration) (tokenManager, error) {
	if signingKey == "" {
		return tokenManager{}, ErrNoSigningKey
	}

	return tokenManager{
		signingKey: signingKey,
		ttl:        ttl,
	}, nil
}

// NewJWT creates new JWT token with claims and subject in payload
func (m tokenManager) NewJWT(sub string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   sub,
		ExpiresAt: time.Now().Add(m.ttl).Unix(),
	})

	return token.SignedString([]byte(m.signingKey))
}

// Parse parses and validating JWT token, returns subject
func (m tokenManager) ParseJWT(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (i interface{}, err error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSignMethod
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", ErrNoClaims
	}

	return claims["sub"].(string), nil
}
