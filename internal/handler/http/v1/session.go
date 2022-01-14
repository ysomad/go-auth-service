package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"

	"github.com/ysomad/go-auth-service/pkg/apperrors"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type sessionHandler struct {
	log logger.Interface
	validation.Gin
	sessionService service.Session
}

func newSessionHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Gin,
	sess service.Session, auth service.Auth) {

	h := &sessionHandler{l, v, sess}

	g := handler.Group("/sessions")
	{
		authenticated := g.Group("/", sessionMiddleware(l, sess))
		{
			secure := authenticated.Group("/", tokenMiddleware(l, auth))
			{
				secure.DELETE(":sessionID", h.terminate)
				secure.DELETE("", h.terminateAll)
			}

			authenticated.GET("", h.get)
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
	currSid, err := sessionID(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminate - sessionID: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := h.sessionService.Terminate(c.Request.Context(), c.Param("sessionID"), currSid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - sessionService - terminate - h.sessionService.Terminate: %w", err))

		if errors.Is(err, apperrors.ErrSessionNotTerminated) {
			abortWithError(c, http.StatusBadRequest, apperrors.ErrSessionNotTerminated)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(c *gin.Context) {
	currSid, err := sessionID(c)
	if err != nil {
		h.log.Error("http - v1 - sessionService - terminateAll - sessionID: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	aid, err := accountID(c)
	if err != nil {
		h.log.Error("http - v1 - sessionService - terminateAll - accountID: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return

	}

	if err := h.sessionService.TerminateAll(c.Request.Context(), aid, currSid); err != nil {
		h.log.Error("http - v1 - sessionService - terminateAll - h.sessionService.TerminateAll: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
