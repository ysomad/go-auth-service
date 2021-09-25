package service

import (
	"context"
	"time"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/pkg/auth"
)

type AuthService struct {
	sessionRepo SessionRepo
	userRepo    UserRepo
	jwtService  auth.JWT
	sessionTTL  time.Duration
}

func NewAuthService(s SessionRepo, u UserRepo, jwt auth.JWT, sessionTTL time.Duration) *AuthService {
	return &AuthService{s, u, jwt, sessionTTL}
}

// ttlSeconds returns TTL in seconds int format
func (as *AuthService) ttlSeconds() int {
	return int(as.sessionTTL.Seconds())
}

// unixExpiration returns unix time dependent on sessionRepo TTL from t
func (as *AuthService) unixExpiration(t time.Time) int64 {
	return t.Add(as.sessionTTL).Unix()
}

// verifyAccess compares current sessionRepo fingerprint, userRepo ip, userRepo agent with received values from client.
// If any of fields are not the same token is considered invalid so userRepo should log in with email and password again
// to receive new refresh token.
func (as *AuthService) verifyAccess(curr *entity.Session, s entity.SessionSecurityDTO) error {
	if curr.ExpiresAt < time.Now().Unix() ||
		curr.Fingerprint != s.Fingerprint ||
		curr.UserIP != s.UserIP ||
		curr.UserAgent != s.UserAgent {
		// TODO: send notification to user when some1 is trying to refresh access token from different location

		return entity.ErrSessionExpired
	}

	return nil
}

// Login identifies userRepo by email, password and creates new refresh sessionRepo.
func (as *AuthService) Login(ctx context.Context, cred entity.UserCredentialsDTO, security entity.SessionSecurityDTO) (entity.JWT, error) {
	// Get userRepo from db
	user, err := as.userRepo.GetByEmail(ctx, cred.Email)
	if err != nil {
		return entity.JWT{}, err
	}

	// Compare passwords
	if err = user.ComparePassword(cred.Password); err != nil {
		return entity.JWT{}, err
	}

	// Generate refresh token
	refreshToken, err := as.jwtService.NewRefresh()
	if err != nil {
		return entity.JWT{}, err
	}

	// Store sessionRepo in redis
	now := time.Now()
	if err = as.sessionRepo.Create(ctx, &entity.Session{
		RefreshToken: refreshToken,
		UserID:       user.ID,
		UserAgent:    security.UserAgent,
		UserIP:       security.UserIP,
		Fingerprint:  security.Fingerprint,
		ExpiresIn:    as.sessionTTL,
		ExpiresAt:    as.unixExpiration(now),
		CreatedAt:    now,
	}); err != nil {
		return entity.JWT{}, err
	}

	// Generate access token
	accessToken, err := as.jwtService.NewAccess(user.ID)
	if err != nil {
		return entity.JWT{}, err
	}

	return entity.JWT{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    as.ttlSeconds(),
	}, nil
}

// RefreshToken creates new access and refresh token pair with expiration time.
// If one of the sessionRepo security fields is not the same as in the current update sessionRepo,
// then the current sessionRepo is deleted and a new one is not created. User should log in with email and password again.
func (as *AuthService) RefreshToken(ctx context.Context, security entity.SessionSecurityDTO) (entity.JWT, error) {
	currSession, err := as.sessionRepo.Get(ctx, security.RefreshToken)
	if err != nil {
		return entity.JWT{}, err
	}

	// Delete current sessionRepo from redis
	if err = as.sessionRepo.Terminate(ctx, currSession.RefreshToken); err != nil {
		return entity.JWT{}, err
	}

	// Check userRepo agent, ip, fingerprint and refresh token lifetime, if it's expire return token expired error
	if err = as.verifyAccess(currSession, security); err != nil {
		return entity.JWT{}, err
	}

	// Generate refresh token
	refreshToken, err := as.jwtService.NewRefresh()
	if err != nil {
		return entity.JWT{}, err
	}

	// Create new sessionRepo in redis
	now := time.Now()
	if err = as.sessionRepo.Create(ctx, &entity.Session{
		RefreshToken: refreshToken,
		UserID:       currSession.UserID,
		UserAgent:    security.UserAgent,
		UserIP:       security.UserIP,
		Fingerprint:  security.Fingerprint,
		ExpiresIn:    as.sessionTTL,
		ExpiresAt:    as.unixExpiration(now),
		CreatedAt:    now,
	}); err != nil {
		return entity.JWT{}, err
	}

	// Generate access token
	accessToken, err := as.jwtService.NewAccess(currSession.UserID)
	if err != nil {
		return entity.JWT{}, err
	}

	return entity.JWT{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    as.ttlSeconds(),
	}, err
}
