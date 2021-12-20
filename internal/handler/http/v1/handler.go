// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

const _path = "/v1"

func SetupHandlers(handler *gin.Engine, v validation.Validator, acc service.Account, s service.Session, a service.Auth) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Swagger UI
	handler.Static(fmt.Sprintf("%s/swagger/", _path), "third_party/swaggerui")

	// Resource handlers
	h := handler.Group(_path)
	{
		newAccountHandler(h, v, acc, s)
		newSessionHandler(h, v, s)
		newAuthHandler(h, v, s, a)
	}
}
