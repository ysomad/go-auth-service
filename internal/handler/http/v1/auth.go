package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/service"

	"github.com/ysomad/go-auth-service/pkg/apperrors"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type authHandler struct {
	log logger.Interface
	validation.Gin
	cfg               *config.Config
	authService       service.Auth
	socialAuthService service.SocialAuth
}

func newAuthHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Gin, cfg *config.Config, s service.Session,
	a service.Auth, sa service.SocialAuth) {

	h := &authHandler{l, v, cfg, a, sa}

	g := handler.Group("/auth")
	{
		g.POST("login", h.login).Use(setCSRFTokenMiddleware(l, cfg))

		social := g.Group("/social", setCSRFTokenMiddleware(l, cfg))
		{
			social.GET("", h.socialAuthorizationURL)
			social.POST("github", h.githubLogin).Use(csrfMiddleware(l, cfg))
		}

		protected := g.Group("/", csrfMiddleware(l, cfg), sessionMiddleware(l, s))
		{
			protected.POST("logout", h.logout)
			protected.POST("token", h.token)

		}
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *authHandler) login(c *gin.Context) {
	var r loginRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		h.log.Info(err.Error())
		abortWithValidationError(c, http.StatusBadRequest, h.TranslateError(err))
		return
	}

	s, err := h.authService.EmailLogin(
		c.Request.Context(),
		r.Email,
		r.Password,
		service.Device{
			IP:        c.ClientIP(),
			UserAgent: c.Request.Header.Get("User-Agent"),
		},
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - login: %w", err))

		if errors.Is(err, apperrors.ErrAccountIncorrectPassword) ||
			errors.Is(err, apperrors.ErrAccountNotFound) {
			abortWithError(c, http.StatusUnauthorized, apperrors.ErrAccountIncorrectEmailOrPassword)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		h.cfg.Session.CookieKey,
		s.ID,
		s.TTL,
		apiPath,
		h.cfg.Session.CookieDomain,
		h.cfg.Session.CookieSecure,
		h.cfg.Session.CookieHTTPOnly,
	)
	c.Status(http.StatusOK)
}

func (h *authHandler) logout(c *gin.Context) {
	sid, err := sessionID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - logout - sessionID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - logout: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		h.cfg.Session.CookieKey,
		"",
		-1,
		apiPath,
		h.cfg.Session.CookieDomain,
		h.cfg.Session.CookieSecure,
		h.cfg.Session.CookieHTTPOnly,
	)
	c.Status(http.StatusNoContent)
}

type tokenRequest struct {
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func (h *authHandler) token(c *gin.Context) {
	var r tokenRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		h.log.Info(err.Error())
		abortWithValidationError(c, http.StatusBadRequest, h.TranslateError(err))
		return
	}

	aid, err := accountID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token - accountID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	t, err := h.authService.NewAccessToken(c.Request.Context(), aid, r.Password)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token: %w", err))

		if errors.Is(err, apperrors.ErrAccountIncorrectPassword) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tokenResponse{t})
}

type getOAuthURIResponse struct {
	URL string `json:"url"`
}

func (h *authHandler) socialAuthorizationURL(c *gin.Context) {
	provider, found := c.GetQuery("provider")
	if !found || provider == "" {
		abortWithError(c, http.StatusBadRequest, apperrors.ErrAuthProviderNotFound)
		return
	}

	uri, err := h.socialAuthService.AuthorizationURL(c.Request.Context(), provider)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - getOAuthURI: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, getOAuthURIResponse{uri.String()})
}

func (h *authHandler) githubLogin(c *gin.Context) {
	code, found := c.GetQuery("code")
	if !found || code == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// TODO: CSRF protection
	state, found := c.GetQuery("state")
	if !found || state == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	s, err := h.socialAuthService.GitHubLogin(
		c.Request.Context(),
		code,
		service.Device{
			UserAgent: c.Request.Header.Get("User-Agent"),
			IP:        c.ClientIP(),
		},
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - githubCallback: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(
		h.cfg.Session.CookieKey,
		s.ID,
		s.TTL,
		apiPath,
		h.cfg.Session.CookieDomain,
		h.cfg.Session.CookieSecure,
		h.cfg.Session.CookieHTTPOnly,
	)
	c.Status(http.StatusOK)
}
