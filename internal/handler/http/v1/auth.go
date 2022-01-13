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

type authHandler struct {
	log logger.Interface
	validation.Gin
	authService       service.Auth
	socialAuthService service.SocialAuth
}

func newAuthHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Gin, s service.Session,
	a service.Auth, sa service.SocialAuth) {

	h := &authHandler{l, v, a, sa}

	g := handler.Group("/auth")
	{
		g.POST("login", h.login)

		social := g.Group("/social")
		{
			social.GET("", h.socialAuthorizationURL)
			social.POST("github", h.githubLogin)
		}

		authenticated := g.Group("/", sessionMiddleware(l, s))
		{
			authenticated.POST("logout", h.logout)
			authenticated.POST("token", h.token)

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
		c.AbortWithStatusJSON(http.StatusBadRequest, h.TranslateError(err))
		return
	}

	cookie, err := h.authService.EmailLogin(
		c.Request.Context(),
		r.Email,
		r.Password,
		service.NewDevice(c.Request.Header.Get("User-Agent"), c.ClientIP()),
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

	setSessionCookie(c, cookie)
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
		c.AbortWithStatusJSON(http.StatusBadRequest, h.TranslateError(err))
		return
	}

	aid, err := accountID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token - accountID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := h.authService.NewAccessToken(c.Request.Context(), aid, r.Password)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token: %w", err))

		if errors.Is(err, apperrors.ErrAccountIncorrectPassword) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tokenResponse{token})
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
	cbErr, found := c.GetQuery("error")
	if found && cbErr != "" {
		desc, _ := c.GetQuery("error_description")
		h.log.Error(fmt.Errorf("http - v1 - auth - githubCallback - %s(%s)", cbErr, desc))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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

	cookie, err := h.socialAuthService.GitHubLogin(
		c.Request.Context(),
		code,
		service.NewDevice(c.Request.Header.Get("User-Agent"), c.ClientIP()),
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - githubCallback: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	setSessionCookie(c, cookie)
	c.Status(http.StatusOK)
}
