package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service/repository"
)

type sessionService struct {
	accountRepo AccountRepo

	sessionRepo SessionRepo
	sessionTTL  time.Duration

	cacheRepo CacheRepo
	cacheTTL  time.Duration
}

func NewSessionService(a AccountRepo, s SessionRepo, c CacheRepo,
	cacheTTL time.Duration, sessionTTL time.Duration) *sessionService {

	return &sessionService{
		accountRepo: a,
		sessionRepo: s,
		cacheRepo:   c,
		sessionTTL:  sessionTTL,
		cacheTTL:    cacheTTL,
	}
}

func (s *sessionService) Create(ctx context.Context, aid string,
	d domain.Device) (domain.Session, error) {

	// TODO: generic errors pkg/httperror

	sess, err := domain.NewSession(aid, d.UserAgent, d.UserIP, s.sessionTTL)
	if err != nil {
		return domain.Session{}, err
	}

	if err = s.sessionRepo.Create(ctx, sess); err != nil {
		return domain.Session{}, err
	}

	return sess, nil
}

func (s *sessionService) Get(ctx context.Context, sid string) (domain.Session, error) {
	// TODO: refactor
	// TODO: make errors generic pkg/httperror

	var sess domain.Session

	err := s.cacheRepo.Get(ctx, repository.BuildCacheKey("ses", sid), &sess)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - s.cache.Get: %w", err)
	}

	if (sess != domain.Session{}) {
		return sess, nil
	}

	sess, err = s.sessionRepo.Get(ctx, sid)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - s.sessionRepo.Get: %w", err)
	}

	if err = s.cacheRepo.Set(ctx, repository.BuildCacheKey("ses", sid), sess, s.cacheTTL); err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Find - s.cache.Set: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetAll(ctx context.Context, aid string) ([]entity.Session, error) {
	// TODO: add caching list of user sessions
	// TODO: refactor

	// TODO: generic errors pkg/httperror
	sessions, err := s.sessionRepo.GetAll(ctx, aid)
	if err != nil {
		return nil, fmt.Errorf("sessionService - s.sessionRepo.GetAll: %w", err)
	}

	return sessions, nil
}

func (s *sessionService) Terminate(ctx context.Context, sid string) error {
	// TODO: generic errors pkg/httperror
	if err := s.cacheRepo.Delete(ctx, repository.BuildCacheKey("ses", sid)); err != nil {
		return fmt.Errorf("sessionService - s.cache.Delete: %w", err)
	}

	if err := s.sessionRepo.Delete(ctx, sid); err != nil {
		return fmt.Errorf("sessionService - s.sessionRepo.Delete: %w", err)
	}

	return nil
}

func (s *sessionService) TerminateAll(ctx context.Context, uid string) error {
	// TODO: add list of sessions to cache
	// TODO: remove list of user sessions from cache first
	// TODO: generic errors pkg/httperror

	if err := s.sessionRepo.DeleteAll(ctx, uid); err != nil {
		return fmt.Errorf("sessionService - s.sessionRepo.DeleteAll: %w", err)
	}

	return nil
}
