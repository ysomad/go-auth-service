package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/pkg/auth"
)

type AuthService struct {
	session   SessionRepo
	user      UserRepo
	jwt       auth.JWT
	expiresIn time.Duration
}

func NewAuthService(s SessionRepo, u UserRepo, m auth.JWT, e time.Duration) *AuthService {
	return &AuthService{s, u, m, e}
}

func (as *AuthService) Login(ctx context.Context, req entity.LoginRequest, s entity.RefreshSession) (entity.LoginResponse, error) {
	// Get user from db
	user, err := as.user.GetByEmail(ctx, req.Email)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	// Compare passwords
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return entity.LoginResponse{}, errors.New(entity.UserIncorrectErr)
	}

	// Generate access and refresh token
	accessToken, err := as.jwt.NewAccess(user.ID)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	refreshToken, err := as.jwt.NewRefresh()
	if err != nil {
		return entity.LoginResponse{}, err
	}

	s.UserID = user.ID
	s.RefreshToken = refreshToken
	s.ExpiresIn = as.expiresIn

	// Create user session in redis
	if err = as.session.Create(s); err != nil {
		return entity.LoginResponse{}, err
	}

	return entity.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
