// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/service"

	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

const apiPath = "/v1"

func SetupHandlers(
	handler *gin.Engine,
	l logger.Interface,
	v validation.Gin,
	cfg *config.Config,
	acc service.Account,
	sess service.Session,
	auth service.Auth,
	social service.SocialAuth,
) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// CORS
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"http://localhost:3000", "https://github.com"}
	corsCfg.AllowCredentials = true

	handler.Use(cors.New(corsCfg))

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Swagger UI
	handler.Static(fmt.Sprintf("%s/swagger/", apiPath), "third_party/swaggerui")

	// Resource handlers
	h := handler.Group(apiPath)
	{
		newAccountHandler(h, l, v, &cfg.Session, acc, sess, auth)
		newSessionHandler(h, l, v, sess, auth)
		newAuthHandler(h, l, v, sess, auth, social)
	}
}
