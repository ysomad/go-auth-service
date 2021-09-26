package v1

import (
	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type userRoutes struct {
	validator   validation.Validator
	userService service.User
}

func newUserRoutes(handler *gin.RouterGroup, t validation.Validator, u service.User, j auth.JWT) {
	r := &userRoutes{t, u}

	h := handler.Group("/users")
	{
		h.POST("", r.create)

		authenticated := h.Group("/", authMiddleware(j))
		{
			authenticated.GET("", r.getUser)
			authenticated.PATCH("", r.partialUpdate)
			authenticated.PATCH("archive", r.archive)
		}
	}
}

type userCreateRequest struct {
	Email           string `json:"email" example:"user@mail.com" binding:"required,email,lte=255"`
	Password        string `json:"password" example:"secret" binding:"required,gte=6,lte=128"`
	ConfirmPassword string `json:"confirmPassword" example:"secret" binding:"required,eqfield=Password"`
}

// @Summary     Create
// @Description Create a new user with email and password
// @ID          userCreate
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body userCreateRequest true "To create a new user email and password should be provided"
// @Success     204
// @Failure     400,500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [post]
func (r *userRoutes) create(c *gin.Context) {
	var req userCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.validator.TranslateAll(err))
		return
	}

	if err := r.userService.Create(c.Request.Context(), req.Email, req.Password); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

type userArchiveRequest struct {
	IsArchive *bool `json:"isArchive" example:"false" binding:"required"`
}

// @Summary     Archive/Restore
// @Description Archive or restore user
// @ID          userArchive
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body userArchiveRequest true "To archive or restore a user is_archive should be provided"
// @Success     204
// @Failure     400,401,500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/archive [patch]
// @Security    Bearer
func (r *userRoutes) archive(c *gin.Context) {
	var req userArchiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.validator.TranslateAll(err))
		return
	}

	id, err := getUserID(c)
	if err != nil {
		abortWithError(c, http.StatusUnauthorized, err)
		return
	}

	if err = r.userService.Archive(c.Request.Context(), id, *req.IsArchive); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

type userPartialUpdateRequest struct {
	Username  string `json:"username" example:"username" binding:"omitempty,alphanum,gte=4,lte=32"`
	FirstName string `json:"firstName" example:"Alex"  binding:"omitempty,alpha,lte=50"`
	LastName  string `json:"lastName" example:"Malykh" binding:"omitempty,alpha,lte=50"`
}

// @Summary     Partial update
// @Description Update user data partially
// @ID         	userPartialUpdate
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body userPartialUpdateRequest true "Provide at least one user field to update user data"
// @Success     204
// @Failure     400,401,500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [patch]
// @Security    Bearer
func (r *userRoutes) partialUpdate(c *gin.Context) {
	var req userPartialUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.validator.TranslateAll(err))
		return
	}

	id, err := getUserID(c)
	if err != nil {
		abortWithError(c, http.StatusUnauthorized, err)
		return
	}

	cols := entity.UpdateColumns{
		"username":   req.Username,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
	}
	if err = cols.Validate(); err != nil {
		abortWithError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if err = r.userService.PartialUpdate(c.Request.Context(), id, cols); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary     Get
// @Description Receive user data
// @ID          userGet
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Failure     400,401,500 {object} messageResponse
// @Success     200 {object} entity.User
// @Router      /users [get]
// @Security    Bearer
func (r *userRoutes) getUser(c *gin.Context) {
	id, err := getUserID(c)
	if err != nil {
		abortWithError(c, http.StatusUnauthorized, err)
		return
	}

	u, err := r.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, u)
}
