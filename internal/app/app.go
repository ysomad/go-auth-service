package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/config"

	v1 "github.com/ysomad/go-auth-service/internal/handler/http/v1"
	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/internal/service/repository"

	"github.com/ysomad/go-auth-service/pkg/httpserver"
	"github.com/ysomad/go-auth-service/pkg/jwt"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/mongodb"
	"github.com/ysomad/go-auth-service/pkg/postgres"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Postgres
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// MongoDB
	mcli, err := mongodb.NewClient(cfg.MongoDB.URI, cfg.MongoDB.Username, cfg.MongoDB.Password)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - mongodb.NewClient: %w", err))
	}
	mdb := mcli.Database(cfg.MongoDB.Database)

	// Redis
	/*
		rdb := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       0,
		})
	*/

	// Service
	accountRepo := repository.NewAccountRepo(pg)
	sessionRepo := repository.NewSessionRepo(mdb)

	accountService := service.NewAccountService(accountRepo)
	sessionService := service.NewSessionService(cfg, sessionRepo)

	jwt, err := jwt.New(cfg.AccessToken.SigningKey, cfg.AccessToken.TTL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - auth.NewTokenManager: %w", err))
	}

	authService := service.NewAuthService(cfg, jwt, accountService, sessionService)
	socialAuthService := service.NewSocialAuthService(cfg, accountService, sessionService)

	// TODO: refactor
	// Validation translator
	v, err := validation.NewGinValidator()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - validation.NewGinValidator: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	v1.SetupHandlers(handler, l, v, cfg, accountService, sessionService, authService, socialAuthService)
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
