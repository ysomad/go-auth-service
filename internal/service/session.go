package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service/repository"
)

const (
	sessionCacheKey = "ses"
)

type sessionService struct {
	userRepo UserRepo

	sessionRepo SessionRepo
	sessionTTL  time.Duration

	cache    CacheRepo
	cacheTTL time.Duration
}

func NewSessionService(u UserRepo, s SessionRepo, c CacheRepo,
	cacheTTL time.Duration, sessionTTL time.Duration) *sessionService {

	return &sessionService{u, s, sessionTTL, c, cacheTTL}
}

func (s *sessionService) LoginWithEmail(ctx context.Context, email, password string,
	d entity.Device) (entity.Session, error) {

	// Get user from DB
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return entity.Session{}, fmt.Errorf("sessionService - LoginWithEmail - userRepo.FindByEmail: %w", err)
	}

	// Compare passwords
	if err = u.ComparePassword(password); err != nil {
		return entity.Session{}, fmt.Errorf("sessionService - LoginWithEmail - u.ComparePassword: %w", err)
	}

	// Create session entity
	sess, err := entity.NewSession(u.ID, d.UserAgent, d.UserIP, s.sessionTTL)
	if err != nil {
		return entity.Session{}, fmt.Errorf("sessionService - LoginWithEmail - entity.NewSession: %w", err)
	}

	// Create session in DB
	if err = s.sessionRepo.Create(ctx, sess); err != nil {
		return entity.Session{}, fmt.Errorf("sessionService - LoginWithEmail - s.sessionRepo.Create: %w", err)
	}

	// Add session to cache
	if err = s.cache.Set(ctx, repository.BuildCacheKey(sessionCacheKey, sess.ID), sess, s.cacheTTL); err != nil {
		return entity.Session{}, fmt.Errorf("sessionService - LoginWithEmail - s.sessionRepo.Create: %w", err)
	}

	return sess, nil
}

func (s *sessionService) Find(ctx context.Context, sid string) (entity.Session, error) {
	var sess entity.Session

	// Find in cache
	err := s.cache.Get(ctx, repository.BuildCacheKey(sessionCacheKey, sid), &sess)
	if err != nil {
		return entity.Session{}, fmt.Errorf("sessionService - s.cache.Get: %w", err)
	}

	if (sess != entity.Session{}) {
		return sess, nil
	}

	// Find in DB
	sess, err = s.sessionRepo.Get(ctx, sid)
	if err != nil {
		return entity.Session{}, fmt.Errorf("sessionService - s.sessionRepo.Get: %w", err)
	}

	// Set session to cache
	if err = s.cache.Set(ctx, repository.BuildCacheKey(sessionCacheKey, sid), sess, s.cacheTTL); err != nil {
		return entity.Session{}, fmt.Errorf("sessionService - Find - s.cache.Set: %w", err)
	}

	return sess, nil
}

func (s *sessionService) FindAll(ctx context.Context, uid string) ([]entity.Session, error) {
	// TODO: add caching list of user sessions

	sessions, err := s.sessionRepo.GetAll(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("sessionService - s.sessionRepo.GetAll: %w", err)
	}

	return sessions, nil
}

func (s *sessionService) Terminate(ctx context.Context, sid string) error {
	if err := s.cache.Delete(ctx, repository.BuildCacheKey(sessionCacheKey, sid)); err != nil {
		return fmt.Errorf("sessionService - s.cache.Delete: %w", err)
	}

	if err := s.sessionRepo.Delete(ctx, sid); err != nil {
		return fmt.Errorf("sessionService - s.sessionRepo.Delete: %w", err)
	}

	return nil
}

func (s *sessionService) TerminateAll(ctx context.Context, uid string) error {
	// TODO: remove list of user sessions from cache first

	if err := s.sessionRepo.DeleteAll(ctx, uid); err != nil {
		return fmt.Errorf("sessionService - s.sessionRepo.DeleteAll: %w", err)
	}

	return nil
}
