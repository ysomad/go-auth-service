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
	validation.Validator
	log         logger.Interface
	authService service.Auth
}

func newAuthHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Validator, s service.Session, a service.Auth) {
	h := &authHandler{v, l, a}

	g := handler.Group("/auth")
	{
		g.POST("login", h.login)

		authenticated := g.Group("/", sessionMiddleware(l, s))
		{
			authenticated.POST("logout", h.logout)
			authenticated.POST("token", h.token)
		}
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=128"`
}

func (h *authHandler) login(c *gin.Context) {
	var r loginRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		h.log.Info(err.Error())
		abortWithValidationError(c, http.StatusBadRequest, h.TranslateError(err))
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

		if errors.Is(err, apperrors.ErrAccountIncorrectPassword) {
			abortWithError(c, http.StatusUnauthorized, apperrors.ErrAccountNotAuthorized)
			return
		}

		if errors.Is(err, apperrors.ErrAccountNotFound) {
			abortWithError(c, http.StatusNotFound, apperrors.ErrAccountNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(domain.SessionCookieKey, cookie.ID(), cookie.TTL(), apiPath, "", false, true) // TODO: set Secure to True
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
	Password string `json:"password" binding:"required,gte=6,lte=128"`
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

	c.JSON(http.StatusOK, token)
}
