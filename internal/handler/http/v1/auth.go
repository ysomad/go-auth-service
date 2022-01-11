package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"

	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type authHandler struct {
	log logger.Interface
	validation.Validator
	authService  service.Auth
	oauthService service.OAuth
}

func newAuthHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Validator, s service.Session,
	a service.Auth, oa service.OAuth) {

	h := &authHandler{l, v, a, oa}

	g := handler.Group("/auth")
	{
		oauth := g.Group("/oauth")
		{
			oauth.GET("", h.getOAuthURI)
			oauth.POST("github", h.githubLogin)
		}

		authenticated := g.Group("/", sessionMiddleware(l, s))
		{
			authenticated.POST("logout", h.logout)
			authenticated.POST("token", h.token)

		}

		g.POST("", h.login)
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *authHandler) setSessionCookie(c *gin.Context, cookie domain.SessionCookie) {
	c.SetCookie(
		cookie.Key,
		cookie.ID,
		cookie.TTL,
		apiPath,
		cookie.Domain,
		cookie.Secure,
		cookie.HTTPOnly,
	)
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
		domain.NewDevice(c.Request.Header.Get("User-Agent"), c.ClientIP()),
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

	h.setSessionCookie(c, cookie)
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
	uri string `json:"uri"`
}

func (h *authHandler) getOAuthURI(c *gin.Context) {
	provider, found := c.GetQuery("provider")
	if !found || provider == "" {
		abortWithError(c, http.StatusBadRequest, apperrors.ErrAuthProviderNotFound)
		return
	}

	uri, err := h.oauthService.GetAuthorizeURI(c.Request.Context(), provider)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - getOAuthURI: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, uri)
	c.JSON(http.StatusOK, getOAuthURIResponse{uri})
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

	cookie, err := h.oauthService.GitHubLogin(
		c.Request.Context(),
		code,
		domain.NewDevice(c.Request.Header.Get("User-Agent"), c.ClientIP()),
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - githubCallback: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.setSessionCookie(c, cookie)
	c.Status(http.StatusOK)
}
