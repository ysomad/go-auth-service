package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type accountHandler struct {
	validation.Validator
	account service.Account
}

func newAccountHandler(handler *gin.RouterGroup, v validation.Validator, u service.Account, s service.Session) {
	r := &accountHandler{v, u}

	h := handler.Group("/users")
	{
		h.POST("", r.create)

		authenticated := h.Group("/", sessionMiddleware(s))
		{
			authenticated.GET("", r.get)
			authenticated.DELETE("", r.archive)
		}
	}
}

type accountCreateRequest struct {
	Email           string `json:"email" example:"account@mail.com" binding:"required,email,lte=255"`
	Password        string `json:"password" example:"secret" binding:"required,gte=8,lte=128"`
	ConfirmPassword string `json:"confirmPassword" example:"secret" binding:"required,eqfield=Password"`
}

func (h *accountHandler) create(c *gin.Context) {
	// TODO: fix
	var r accountCreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, h.TranslateError(err))
		return
	}

	acc, err := h.account.Create(c.Request.Context(), r.Email, r.Password)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, acc)
}

func (h *accountHandler) archive(c *gin.Context) {
	// TODO: fix
	// get account id from context
	var aid string

	if err := h.account.Archive(c.Request.Context(), aid); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (r *accountHandler) get(c *gin.Context) {
	// TODO: fix
	// get account id from context
	var aid string

	acc, err := r.account.GetByID(c.Request.Context(), aid)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, acc)
}
