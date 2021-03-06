package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/service"

	"github.com/ysomad/go-auth-service/pkg/apperrors"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/utils"
)

func tokenMiddleware(l logger.Interface, a service.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		aid, err := accountID(c)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - accountID: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		t, found := c.GetQuery("token")
		if !found || t == "" {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - c.GetQuery: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		sub, err := a.ParseAccessToken(c.Request.Context(), t)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - auth.ParseAccessToken: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if sub != aid {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func sessionMiddleware(l logger.Interface, s service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("id")
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - c.Cookie: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		s, err := s.GetByID(c.Request.Context(), sid)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - s.Get: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		d := service.Device{UserAgent: c.Request.Header.Get("User-Agent"), IP: c.ClientIP()}

		if s.IP != d.IP || s.UserAgent != d.UserAgent {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware: %w", apperrors.ErrSessionDeviceMismatch))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("sid", s.ID)
		c.Set("aid", s.AccountID)
		c.Next()
	}
}

func csrfMiddleware(l logger.Interface, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pt string

		qt := c.Query("state")
		ht := c.Request.Header.Get(cfg.CSRFToken.HeaderKey)
		if qt != "" {
			pt = qt
		} else if ht != "" {
			pt = ht
		} else {
			l.Error(fmt.Errorf("http - v1 - middleware - csrfMiddleware - c.Request.Header.Get: %w", apperrors.ErrCSRFTokenPublicNotFound))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		ct, err := c.Cookie(cfg.CSRFToken.CookieKey)
		if err != nil || ct == "" {
			l.Error(fmt.Errorf("http - v1 - middleware - csrfMiddleware - c.Cookie: %w", apperrors.ErrCSRFTokenCookieNotFound))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if pt != ct {
			l.Error(fmt.Errorf("http - v1 - middleware - csrfMiddleware: %w", apperrors.ErrCSRFDetected))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		t, err := utils.UniqueString(32)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - csrfMiddleware: %w", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Next()

		c.Header(cfg.CSRFToken.HeaderKey, t)
		c.SetCookie(
			cfg.CSRFToken.CookieKey,
			t,
			int(cfg.CSRFToken.TTL.Seconds()),
			apiPath,
			"",
			cfg.Session.CookieSecure,
			cfg.Session.CookieHTTPOnly,
		)
	}
}

func setCSRFTokenMiddleware(l logger.Interface, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := utils.UniqueString(32)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - setCSRFTokenMiddleware: %w", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Next()

		c.Header(cfg.CSRFToken.HeaderKey, t)
		c.SetCookie(
			cfg.CSRFToken.CookieKey,
			t,
			int(cfg.CSRFToken.TTL.Seconds()),
			apiPath,
			"",
			cfg.Session.CookieSecure,
			cfg.Session.CookieHTTPOnly,
		)
	}
}

// accountID returns account id from context
func accountID(c *gin.Context) (string, error) {
	aid := c.GetString("aid")

	_, err := uuid.Parse(aid)
	if err != nil {
		return "", apperrors.ErrAccountContextNotFound
	}

	return aid, nil
}

// sessionID return session id from context
func sessionID(c *gin.Context) (string, error) {
	sid := c.GetString("sid")

	if sid == "" {
		return "", apperrors.ErrSessionContextNotFound
	}

	return sid, nil
}
