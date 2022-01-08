package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"

	"github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/logger"
)

func tokenMiddleware(log logger.Interface, authService service.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		aid, err := accountID(c)
		if err != nil {
			log.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - accountID: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, found := c.GetQuery("token")
		if !found || token == "" {
			log.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - c.GetQuery: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		sub, err := authService.ParseAccessToken(c.Request.Context(), token)
		if err != nil {
			log.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - authService.ParseAccessToken: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if sub != aid {
			log.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func sessionMiddleware(log logger.Interface, sessionService service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("id")
		if err != nil {
			log.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - c.Cookie: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, err := sessionService.Get(c.Request.Context(), sid)
		if err != nil {
			log.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - s.Get: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		device := domain.NewDevice(c.Request.UserAgent(), c.ClientIP())

		if session.IP != device.IP || session.UserAgent != device.UserAgent {
			log.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware: %w", errors.ErrSessionMismatchedDevice))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("sid", session.ID)
		c.Set("aid", session.AccountID)
		c.Next()
	}
}

func accountID(c *gin.Context) (string, error) {
	aid := c.GetString("aid")

	_, err := uuid.Parse(aid)
	if err != nil {
		return "", errors.ErrAccountNotInContext
	}

	return aid, nil
}

func sessionID(c *gin.Context) (string, error) {
	sid := c.GetString("sid")

	if sid == "" {
		return "", errors.ErrSessionNotInContext
	}

	return sid, nil
}
