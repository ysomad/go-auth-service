// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

func SetupHandlers(handler *gin.Engine, v validation.Validator, a service.Account, s service.Session) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Resource handlers
	h := handler.Group("/v1")
	{
		newAccountHandler(h, v, a, s)
		newSessionHandler(h, v, s)
		newAuthHandler(h, v, s)
	}
}
