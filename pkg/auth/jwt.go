package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt"
)

type JWT interface {
	NewAccess(userID uuid.UUID) (string, error)
	NewRefresh() (uuid.UUID, error)
	Validate(accessToken string) (jwt.MapClaims, error)
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

func (m *JWTManager) NewAccess(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userID.String(),
		ExpiresAt: time.Now().Add(m.accessTokenTTL).Unix(),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *JWTManager) NewRefresh() (uuid.UUID, error) {
	return uuid.NewRandom()
}

// Validate parses and validating JWT token, returns user id from it
func (m *JWTManager) Validate(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("error get user claims from token")
	}

	return claims, nil
}
