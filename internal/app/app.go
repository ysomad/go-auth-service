// Package app configures and runs application.
package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/config"
	v1 "github.com/ysomad/go-auth-service/internal/delivery/http/v1"
	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/internal/service/repo"
	"github.com/ysomad/go-auth-service/pkg/httpserver"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		logger.Fatal(err, "app - Run - postgres.NewPostgres")
	}
	defer pg.Close()

	// Service
	userService := service.NewUserService(repo.NewUserRepo(pg))

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, userService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Error(err, "app - Run - httpServer.Notify")
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(err, "app - Run - httpServer.Shutdown")
	}
}
