package v1

import (
	"context"
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
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		sub, err := authService.ParseAccessToken(context.Background(), token)
		if err != nil {
			log.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - authService.ParseAccessToken: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// sub - account id from token payload, aid - account id from context
		if sub != aid {
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

		ctx := c.Request.Context()

		sess, err := sessionService.Get(ctx, sid)
		if err != nil {
			log.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - s.Get: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		d := domain.NewDevice(c.Request.Header.Get("User-Agent"), c.ClientIP())

		// Check current request device vs session device
		if sess.IP != d.IP || sess.UserAgent != d.UserAgent {
			// TODO: send notification that someone logged in on new device
			sessionService.Terminate(ctx, sid)

			c.Status(http.StatusUnauthorized)
			return
		}

		c.Set("sid", sess.ID)
		c.Set("aid", sess.AccountID)
		c.Next()
	}
}

func accountID(c *gin.Context) (string, error) {
	aid := c.GetString("aid")

	_, err := uuid.Parse(aid)
	if err != nil {
		return "", fmt.Errorf("uuid.Parse: %w", err)
	}

	return aid, nil
}

func sessionID(c *gin.Context) (string, error) {
	sid := c.GetString("sid")

	if sid == "" {
		return "", errors.ErrSessionNotFound
	}

	return sid, nil
}
