package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type sessionHandler struct {
	validation.Validator
	session service.Session
}

func newSessionHandler(handler *gin.RouterGroup, v validation.Validator, s service.Session) {
	h := &sessionHandler{v, s}

	g := handler.Group("/auth")
	{
		authenticated := g.Group("/", sessionMiddleware(s))
		{
			authenticated.GET("", h.get)
			authenticated.DELETE("", h.terminate)
			authenticated.DELETE("", h.terminateAll)
		}
	}
}

func (h *sessionHandler) get(c *gin.Context) {
	panic("implement")

	var sessions []string

	c.JSON(http.StatusOK, sessions)
}

func (h *sessionHandler) terminate(c *gin.Context) {
	panic("implement")

	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(c *gin.Context) {
	panic("implement")

	c.Status(http.StatusNoContent)
}
