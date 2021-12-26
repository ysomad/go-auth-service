package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"

	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/logger"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type accountHandler struct {
	log logger.Interface
	validation.Validator
	accountService service.Account
}

func newAccountHandler(handler *gin.RouterGroup, l logger.Interface, v validation.Validator, u service.Account, s service.Session) {
	h := &accountHandler{l, v, u}

	g := handler.Group("/accounts")
	{
		g.POST("", h.create)

		authenticated := g.Group("/", sessionMiddleware(l, s))
		{
			authenticated.GET("", h.get)
			authenticated.DELETE("", h.archive)
		}
	}
}

type accountCreateRequest struct {
	Email           string `json:"email" binding:"required,email,lte=255"`
	Password        string `json:"password" binding:"required,gte=8,lte=128"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

func (h *accountHandler) create(c *gin.Context) {
	var r accountCreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithValidationError(c, http.StatusBadRequest, h.TranslateError(err))
		return
	}

	if err := h.accountService.Create(c.Request.Context(), r.Email, r.Password); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - create: %w", err))

		if errors.Is(err, apperrors.ErrAccountAlreadyExist) {
			abortWithError(c, http.StatusConflict, apperrors.ErrAccountAlreadyExist)
			return
		}

		if errors.Is(err, apperrors.ErrAccountPasswordNotGenerated) {
			c.AbortWithStatus(http.StatusInternalServerError)
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

	if err := h.accountService.Archive(c.Request.Context(), aid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - archive: %w", err))

		if errors.Is(err, apperrors.ErrAccountNotFound) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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

		if errors.Is(err, apperrors.ErrAccountNotFound) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	acc.Sanitize()

	c.JSON(http.StatusOK, acc)
}
