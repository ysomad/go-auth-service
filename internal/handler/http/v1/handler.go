// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"

	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

const apiPath = "/v1"

func SetupHandlers(
	handler *gin.Engine,
	l logger.Interface,
	v validation.Validator,
	acc service.Account,
	sess service.Session,
	auth service.Auth,
) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Swagger UI
	handler.Static(fmt.Sprintf("%s/swagger/", apiPath), "third_party/swaggerui")

	// Resource handlers
	h := handler.Group(apiPath)
	{
		newAccountHandler(h, l, v, acc, sess, auth)
		newSessionHandler(h, l, v, sess)
		newAuthHandler(h, l, v, sess, auth)
	}
}
