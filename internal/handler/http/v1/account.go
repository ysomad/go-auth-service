package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"

	"github.com/ysomad/go-auth-service/pkg/apperrors"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type accountHandler struct {
	log logger.Interface
	validation.Gin
	cfg            *config.Config
	accountService service.Account
	sessionService service.Session
}

func newAccountHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Gin, cfg *config.Config,
	acc service.Account, s service.Session, auth service.Auth) {

	h := &accountHandler{l, v, cfg, acc, s}

	g := handler.Group("/accounts")
	{
		g.POST("", h.create)

		authenticated := g.Group("/", sessionMiddleware(l, s))
		{
			authenticated.GET("", h.get)

			secure := authenticated.Group("/", tokenMiddleware(l, auth))
			{
				secure.DELETE("", h.archive)
			}
		}
	}
}

type accountCreateRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Username string `json:"username" binding:"required,alphanum,gte=4,lte=16"`
	Password string `json:"password" binding:"required,gte=8,lte=64"`
}

func (h *accountHandler) create(c *gin.Context) {
	var r accountCreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, h.TranslateError(err))
		return
	}

	_, err := h.accountService.Create(
		c.Request.Context(),
		domain.Account{Email: r.Email, Username: r.Username, Password: r.Password},
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - create: %w", err))

		if errors.Is(err, apperrors.ErrAccountAlreadyExist) {
			abortWithError(c, http.StatusConflict, apperrors.ErrAccountAlreadyExist)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *accountHandler) archive(c *gin.Context) {
	aid, err := accountID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - archive - accountID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sid, err := sessionID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - archive - sessionID: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.accountService.Delete(c.Request.Context(), aid, sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - archive - h.accountService.Delete: %w", err))
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

func (h *accountHandler) get(c *gin.Context) {
	aid, err := accountID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - archive - accountID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	acc, err := h.accountService.GetByID(c.Request.Context(), aid)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - get: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, acc)
}
