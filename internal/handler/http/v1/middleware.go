package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ysomad/go-auth-service/internal/service"

	"github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/logger"
)

func tokenMiddleware(log logger.Interface, auth service.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		aid, err := accountID(c)
		if err != nil {
			log.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - accountID: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		t, found := c.GetQuery("token")
		if !found || t == "" {
			log.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - c.GetQuery: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		sub, err := auth.ParseAccessToken(c.Request.Context(), t)
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

func sessionMiddleware(log logger.Interface, session service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("id")
		if err != nil {
			log.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - c.Cookie: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		s, err := session.GetByID(c.Request.Context(), sid)
		if err != nil {
			log.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - s.Get: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		d := service.NewDevice(c.Request.UserAgent(), c.ClientIP())

		if s.IP != d.IP || s.UserAgent != d.UserAgent {
			log.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware: %w", errors.ErrSessionMismatchedDevice))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("sid", s.ID)
		c.Set("aid", s.AccountID)
		c.Next()
	}
}

// accountID returns account id from context
func accountID(c *gin.Context) (string, error) {
	aid := c.GetString("aid")

	_, err := uuid.Parse(aid)
	if err != nil {
		return "", errors.ErrAccountContextNotFound
	}

	return aid, nil
}

// sessionID return session id from context
func sessionID(c *gin.Context) (string, error) {
	sid := c.GetString("sid")

	if sid == "" {
		return "", errors.ErrSessionContextNotFound
	}

	return sid, nil
}
