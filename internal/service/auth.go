package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/pkg/auth"
)

type authService struct {
	sessionRepo SessionRepo
	userRepo    UserRepo
	jwtService  auth.JWT
	sessionTTL  time.Duration
}

func NewAuthService(s SessionRepo, u UserRepo, jwt auth.JWT, sessionTTL time.Duration) *authService {
	return &authService{s, u, jwt, sessionTTL}
}

// ttlSeconds returns TTL in seconds int format
func (s *authService) ttlSeconds() int {
	return int(s.sessionTTL.Seconds())
}

// unixExpiration returns unix time dependent on sessionTTL from t
func (s *authService) unixExpiration(t time.Time) int64 {
	return t.Add(s.sessionTTL).Unix()
}

// verifyAccess compares current session fingerprint, user ip, user agent with received values from client.
// If any of fields are not the same, refresh token is considered invalid. User should log in with email and password
// to receive new refresh token
func (s *authService) verifyAccess(curr *entity.Session, security entity.SessionSecurityDTO) error {
	if curr.ExpiresAt < time.Now().Unix() ||
		curr.Fingerprint != security.Fingerprint ||
		curr.UserIP != security.UserIP ||
		curr.UserAgent != security.UserAgent {
		// TODO: send notification to user when some1 is trying to refresh access token from different location

		return entity.ErrSessionExpired
	}

	return nil
}

// Login identifies user by email, password and creates new refresh session
func (s *authService) Login(ctx context.Context, cred entity.UserCredentialsDTO, security entity.SessionSecurityDTO) (entity.JWT, error) {
	// GetOne user from db
	user, err := s.userRepo.GetByEmail(ctx, cred.Email)
	if err != nil {
		return entity.JWT{}, fmt.Errorf("authService - Login - s.userRepo.GetByEmail: %w", err)
	}

	// Compare passwords
	if err = user.ComparePassword(cred.Password); err != nil {
		return entity.JWT{}, fmt.Errorf("authService - Login - user.ComparePassword: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.jwtService.NewRefresh()
	if err != nil {
		return entity.JWT{}, fmt.Errorf("authService - Login - s.jwtService.NewRefresh: %w", err)
	}

	// Create refresh session in redis
	now := time.Now()
	if err = s.sessionRepo.Create(ctx, &entity.Session{
		RefreshToken: refreshToken,
		UserID:       user.ID,
		UserAgent:    security.UserAgent,
		UserIP:       security.UserIP,
		Fingerprint:  security.Fingerprint,
		ExpiresIn:    s.sessionTTL,
		ExpiresAt:    s.unixExpiration(now),
		CreatedAt:    now,
	}); err != nil {
		return entity.JWT{}, fmt.Errorf("authService - Login - s.sessionRepo.Create: %w", err)
	}

	// Generate access token
	accessToken, err := s.jwtService.NewAccess(user.ID)
	if err != nil {
		return entity.JWT{}, fmt.Errorf("authService - Login - s.jwtService.NewAccess: %w", err)
	}

	return entity.JWT{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.ttlSeconds(),
	}, nil
}

// RefreshToken creates new access and refresh token pair with expiration time.
// If one of session security fields is not the same as in the current refresh session,
// then the current session is deleted and a new one is not created. User should log in with email and password again
func (s *authService) RefreshToken(ctx context.Context, security entity.SessionSecurityDTO) (entity.JWT, error) {
	currSession, err := s.sessionRepo.GetOne(ctx, security.RefreshToken)
	if err != nil {
		return entity.JWT{}, fmt.Errorf("authService - RefreshToken - s.sessionRepo.GetOne: %w", err)
	}

	// Delete current refresh session from redis
	if err = s.sessionRepo.Terminate(ctx, currSession.RefreshToken); err != nil {
		return entity.JWT{}, fmt.Errorf("authService - RefreshToken - s.sessionRepo.Terminate: %w", err)
	}

	// Check user agent, ip, fingerprint and refresh token lifetime
	if err = s.verifyAccess(currSession, security); err != nil {
		return entity.JWT{}, fmt.Errorf("authService - RefreshToken - s.verifyAccess: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.jwtService.NewRefresh()
	if err != nil {
		return entity.JWT{}, fmt.Errorf("authService - RefreshToken - s.jwtService.NewRefresh: %w", err)
	}

	// Create new refresh session in redis
	now := time.Now()
	if err = s.sessionRepo.Create(ctx, &entity.Session{
		RefreshToken: refreshToken,
		UserID:       currSession.UserID,
		UserAgent:    security.UserAgent,
		UserIP:       security.UserIP,
		Fingerprint:  security.Fingerprint,
		ExpiresIn:    s.sessionTTL,
		ExpiresAt:    s.unixExpiration(now),
		CreatedAt:    now,
	}); err != nil {
		return entity.JWT{}, fmt.Errorf("authService - RefreshToken - s.sessionRepo.Create: %w", err)
	}

	// Generate access token
	accessToken, err := s.jwtService.NewAccess(currSession.UserID)
	if err != nil {
		return entity.JWT{}, fmt.Errorf("authService - RefreshToken - s.jwtService.NewAccess: %w", err)
	}

	return entity.JWT{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.ttlSeconds(),
	}, err
}
