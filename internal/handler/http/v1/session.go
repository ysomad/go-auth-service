package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"

	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type sessionHandler struct {
	log logger.Interface
	validation.Validator
	sessionService service.Session
}

func newSessionHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Validator,
	sess service.Session, auth service.Auth) {

	h := &sessionHandler{l, v, sess}

	g := handler.Group("/sessions")
	{
		authenticated := g.Group("/", sessionMiddleware(l, sess))
		{
			authenticated.GET("", h.get)

			secure := authenticated.Group("/", tokenMiddleware(l, auth))
			secure.DELETE(":sessionID", h.terminate)
			secure.DELETE("", h.terminateAll)
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
	currentSid, err := sessionID(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminate - sessionID: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sid := c.Param("sessionID")

	if currentSid == sid {
		abortWithError(c, http.StatusBadRequest, apperrors.ErrSessionNotTerminated)
		return
	}

	if err := h.sessionService.Terminate(c.Request.Context(), sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - terminate - h.sessionService.Terminate: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(c *gin.Context) {
	currentSid, err := sessionID(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - sessionID: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	aid, err := accountID(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - accountID: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return

	}

	if err := h.sessionService.TerminateAll(c.Request.Context(), aid, currentSid); err != nil {
		h.log.Error("http - v1 - session - terminateAll: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
