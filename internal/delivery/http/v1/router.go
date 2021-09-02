// Package v1 implements routing paths. Each services in own file.
package v1

import (
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
// @description REST API authentication service
// @version     1.0
// @host        0.0.0.0:8080
// @BasePath    /api/v1

func NewRouter(handler *gin.Engine, userService service.User) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Routers
	h := handler.Group("/api/v1")
	{
		newUserRoutes(h, userService)
	}
}
