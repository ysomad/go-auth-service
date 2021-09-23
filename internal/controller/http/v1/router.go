// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/ysomad/go-auth-service/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/ysomad/go-auth-service/docs"
	"github.com/ysomad/go-auth-service/internal/service"
)

// Swagger spec:
// @title       Golang auth service
// @description REST API authentication and user management service
// @version     1.0
// @host        0.0.0.0:8080
// @BasePath    /v1

func NewRouter(handler *gin.Engine, t validation.Translator, u service.User, a service.Auth) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("/v1")
	{
		newUserRoutes(h, t, u)
		newAuthRoutes(h, t, a)
	}
}
