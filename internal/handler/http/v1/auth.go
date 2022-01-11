package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/github"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"

	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/util"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type authHandler struct {
	log logger.Interface
	validation.Validator
	sessionCfg  config.Session
	authCfg     config.Auth
	authService service.Auth
}

func newAuthHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Validator,
	cfg config.Config, s service.Session, a service.Auth) {

	h := &authHandler{l, v, cfg.Session, cfg.Auth, a}

	g := handler.Group("/auth")
	{
		login := g.Group("/login")
		{
			gh := login.Group("/github")
			{
				gh.GET("callback", h.githubCallback)
				gh.GET("", h.githubLogin)
			}

			login.POST("", h.login)
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

func (h *authHandler) setSessionCookie(c *gin.Context, cookie domain.SessionCookie) {
	c.SetCookie(
		h.sessionCfg.CookieKey,
		cookie.ID(),
		cookie.TTL(),
		apiPath,
		h.sessionCfg.CookieDomain,
		h.sessionCfg.CookieSecure,
		h.sessionCfg.CookieHttpOnly,
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

func (h *authHandler) githubLogin(c *gin.Context) {
	s, err := util.UniqueString(32)
	if err != nil {
		h.log.Error("http - v1 - auth - githubLogin - util.UniqueString: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	url, err := url.Parse(github.Endpoint.AuthURL)
	if err != nil {
		h.log.Error("http - v1 - auth - githubLogin - url.Parse: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	q := url.Query()
	q.Set("client_id", h.authCfg.GitHubClientID)
	q.Set("scope", h.authCfg.GitHubScope)
	q.Set("state", s)
	url.RawQuery = q.Encode()

	h.log.Error(url.String())

	c.Redirect(http.StatusTemporaryRedirect, url.String())
}

func (h *authHandler) githubCallback(c *gin.Context) {
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

	cookie, err := h.authService.GitHubLogin(
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
