package jwt

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

type Token interface {
	New(sub string) (string, error)
	Parse(token string) (string, error)
}

type jwtToken struct {
	signingKey string
	ttl        time.Duration
}

func New(signingKey string, ttl time.Duration) (jwtToken, error) {
	if signingKey == "" {
		return jwtToken{}, ErrNoSigningKey
	}

	return jwtToken{
		signingKey: signingKey,
		ttl:        ttl,
	}, nil
}

// New creates new JWT token with claims and subject in payload
func (m jwtToken) New(sub string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   sub,
		ExpiresAt: time.Now().Add(m.ttl).Unix(),
	})

	return token.SignedString([]byte(m.signingKey))
}

// Parse parses and validating JWT token, returns subject
func (m jwtToken) Parse(token string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (i interface{}, err error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSignMethod
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok && !t.Valid {
		return "", ErrNoClaims
	}

	return claims["sub"].(string), nil
}
