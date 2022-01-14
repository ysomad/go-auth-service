// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"fmt"
	"strings"

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
	corsCfg.AllowOrigins = strings.Split(cfg.HTTP.CORSAllowOrigins, " ")
	corsCfg.AllowCredentials = true

	handler.Use(cors.New(corsCfg))

	// Swagger UI
	handler.Static(fmt.Sprintf("%s/swagger/", apiPath), "third_party/swaggerui")

	// Resource handlers
	h := handler.Group(apiPath)
	{
		newAccountHandler(h, l, v, cfg, acc, sess, auth)
		newSessionHandler(h, l, v, sess, auth)
		newAuthHandler(h, l, v, cfg, sess, auth, social)
	}
}
