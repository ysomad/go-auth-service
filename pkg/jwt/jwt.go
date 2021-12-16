package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ErrEmptySigningKey         = errors.New("empty signing key")
	ErrClaims                  = errors.New("error getting claims from token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

type JWT interface {
	New(sub uuid.UUID) (string, error)
	Validate(jwt string) (string, error)
}

type jwtManager struct {
	signingKey string
	ttl        time.Duration
}

func NewJWTManager(signingKey string, ttl time.Duration) (jwtManager, error) {
	if signingKey == "" {
		return jwtManager{}, ErrEmptySigningKey
	}

	return jwtManager{
		signingKey: signingKey,
		ttl:        ttl,
	}, nil
}

// New creates new JWT token with claims and subject
func (m jwtManager) New(sub uuid.UUID) (string, error) {
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.StandardClaims{
		Subject:   sub.String(),
		ExpiresAt: time.Now().Add(m.ttl).Unix(),
	})

	return token.SignedString([]byte(m.signingKey))
}

// Validate parses and validating JWT token, returns user id from it
func (m jwtManager) Validate(jwt string) (string, error) {
	token, err := jwtlib.Parse(jwt, func(token *jwtlib.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok && !token.Valid {
		return "", ErrClaims
	}

	return claims["sub"].(string), nil
}
