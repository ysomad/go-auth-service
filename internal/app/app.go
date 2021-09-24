// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ysomad/go-auth-service/pkg/auth"
	"github.com/ysomad/go-auth-service/pkg/validation"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/config"
	v1 "github.com/ysomad/go-auth-service/internal/controller/http/v1"
	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/internal/service/repo"
	"github.com/ysomad/go-auth-service/pkg/httpserver"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository

	// Postgres
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	// Service
	userRepo := repo.NewUserRepo(pg)

	jwtManager, err := auth.NewJWTManager(cfg.JWT.SigningKey, cfg.JWT.AccessTokenTTL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - auth.NewTokenManager: %w", err))
	}

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(repo.NewSessionRepo(rdb), userRepo, jwtManager, cfg.JWT.RefreshTokenTTL)

	// Validation translator
	trans, err := validation.NewGinValidator()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - validation.NewGinValidator: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, trans, userService, authService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
