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

	g := handler.Group("/sessions")
	{
		authenticated := g.Group("/", sessionMiddleware(s))
		{
			authenticated.DELETE(":sessionID", h.terminate)
			authenticated.GET("", h.get)
			authenticated.DELETE("", h.terminateAll)
		}
	}
}

func (h *sessionHandler) get(c *gin.Context) {
	var sessions []string

	c.JSON(http.StatusOK, sessions)
}

func (h *sessionHandler) terminate(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
