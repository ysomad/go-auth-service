package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/pkg/auth"
)

type AuthService struct {
	session    SessionRepo
	user       UserRepo
	jwt        auth.JWT
	sessionTTL time.Duration
}

func NewAuthService(s SessionRepo, u UserRepo, m auth.JWT, e time.Duration) *AuthService {
	return &AuthService{s, u, m, e}
}

// ttlSeconds returns TTL in seconds int format
func (as *AuthService) ttlSeconds() int {
	return int(as.sessionTTL.Seconds())
}

// unixExpiration returns unix time dependent on session TTL from t
func (as *AuthService) unixExpiration(t time.Time) int64 {
	return t.Add(as.sessionTTL).Unix()
}

// allowAccess compares current session fingerprint, user ip, user agent with received values from client.
// If any of fields are not the same token is considered invalid so user should log in with email and password again
// to receive new refresh token.
func (as *AuthService) allowAccess(curr entity.Session, dto entity.SessionSecurityDTO) bool {
	if curr.ExpiresAt < time.Now().Unix() || curr.Fingerprint != dto.Fingerprint || curr.UserIP != dto.UserIP || curr.UserAgent != dto.UserAgent {
		// TODO: send notification to user when some1 is trying to refresh access token from different location

		return false
	}

	return true
}

func (as *AuthService) Login(ctx context.Context, req entity.LoginRequest, dto entity.SessionSecurityDTO) (entity.LoginResponse, error) {
	// Get user from db
	u, err := as.user.GetByEmail(ctx, req.Email)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	// Compare passwords
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return entity.LoginResponse{}, entity.UserIncorrectErr
	}

	// Generate refresh token
	r, err := as.jwt.NewRefresh()
	if err != nil {
		return entity.LoginResponse{}, err
	}

	// Store session in redis
	now := time.Now()
	if err = as.session.Create(ctx, entity.Session{
		RefreshToken: r,
		UserID:       u.ID,
		UserAgent:    dto.UserAgent,
		UserIP:       dto.UserIP,
		Fingerprint:  dto.Fingerprint,
		ExpiresIn:    as.sessionTTL,
		ExpiresAt:    as.unixExpiration(now),
		CreatedAt:    now,
	}); err != nil {
		return entity.LoginResponse{}, err
	}

	// Generate access token
	a, err := as.jwt.NewAccess(u.ID)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	return entity.LoginResponse{
		AccessToken:  a,
		RefreshToken: r,
		ExpiresIn:    as.ttlSeconds(),
	}, nil
}

func (as *AuthService) RefreshToken(ctx context.Context, dto entity.SessionSecurityDTO) (entity.LoginResponse, error) {
	s, err := as.session.Get(ctx, dto.RefreshToken)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	// Delete current session from redis
	if err = as.session.Terminate(ctx, dto.RefreshToken); err != nil {
		return entity.LoginResponse{}, err
	}

	// Check user agent, ip, fingerprint and refresh token lifetime, if it's expire return token expired error
	if !as.allowAccess(s, dto) {
		return entity.LoginResponse{}, entity.TokenExpiredErr
	}

	// Generate refresh token
	r, err := as.jwt.NewRefresh()
	if err != nil {
		return entity.LoginResponse{}, err
	}

	// Create new session in redis
	now := time.Now()
	if err = as.session.Create(ctx, entity.Session{
		RefreshToken: r,
		UserID:       s.UserID,
		UserAgent:    s.UserAgent,
		UserIP:       s.UserIP,
		Fingerprint:  s.Fingerprint,
		ExpiresIn:    as.sessionTTL,
		ExpiresAt:    as.unixExpiration(now),
		CreatedAt:    now,
	}); err != nil {
		return entity.LoginResponse{}, err
	}

	// Generate access token
	a, err := as.jwt.NewAccess(s.UserID)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	return entity.LoginResponse{
		AccessToken:  a,
		RefreshToken: r,
		ExpiresIn:    as.ttlSeconds(),
	}, err
}
