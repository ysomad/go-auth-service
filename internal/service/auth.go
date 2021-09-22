package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/pkg/auth"
)

type AuthService struct {
	session          SessionRepo
	user             UserRepo
	jwt              auth.JWT
	sessionExpiresIn time.Duration
}

func NewAuthService(s SessionRepo, u UserRepo, m auth.JWT, e time.Duration) *AuthService {
	return &AuthService{s, u, m, e}
}

func (as *AuthService) Login(ctx context.Context, req entity.LoginRequest, s entity.RefreshSession) (entity.LoginResponse, error) {
	// Get user from db
	u, err := as.user.GetByEmail(ctx, req.Email)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	// Compare passwords
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return entity.LoginResponse{}, entity.UserIncorrectCredsErr
	}

	// Generate access and refresh token
	a, err := as.jwt.NewAccess(u.ID)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	r, err := as.jwt.NewRefresh()
	if err != nil {
		return entity.LoginResponse{}, err
	}

	// Set refresh session public fields
	s.UserID = u.ID
	s.RefreshToken = r
	s.ExpiresIn = as.sessionExpiresIn
	s.CreatedAt = time.Now()

	// Create user session in redis
	if err = as.session.Create(s); err != nil {
		return entity.LoginResponse{}, err
	}

	return entity.LoginResponse{
		AccessToken:  a,
		RefreshToken: r,
	}, nil
}
