package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"

	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type sessionHandler struct {
	log logger.Interface
	validation.Validator
	sessionService service.Session
}

func newSessionHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Validator, s service.Session) {
	h := &sessionHandler{l, v, s}

	g := handler.Group("/sessions")
	{
		authenticated := g.Group("/", sessionMiddleware(l, s))
		{
			authenticated.DELETE(":sessionID", h.terminate)
			authenticated.GET("", h.get)
			authenticated.DELETE("", h.terminateAll)
		}
	}
}

func (h *sessionHandler) get(c *gin.Context) {
	aid, err := accountID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - get - accountID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessions, err := h.sessionService.GetAll(c.Request.Context(), aid)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - get: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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
