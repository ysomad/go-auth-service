package v1

import (
	"github.com/ysomad/go-auth-service/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/entity"
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
		h.POST("", r.signUp)
		h.PATCH("archive", r.archive).Use(authMiddleware(j))
		h.GET("", r.getUser).Use(authMiddleware(j))
		h.PATCH("", r.partialUpdate).Use(authMiddleware(j))
	}
}

// @Summary     Create new user
// @Description Create a new user with email and password
// @ID          signup
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body entity.CreateUserRequest true "To create a new user email and password should be provided"
// @Success     204
// @Failure     400 {object} messageResponse
// @Failure     500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [post]
func (r *userRoutes) signUp(c *gin.Context) {
	var req entity.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.validator.TranslateAll(err))
		return
	}

	if err := r.userService.SignUp(c.Request.Context(), req); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary     Archive or restore user
// @Description Archive or restore user
// @ID          archive
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body entity.ArchiveUserRequest true "To archive or restore a user is_archive should be provided"
// @Success     204
// @Failure     401 {object} messageResponse
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Failure     500 {object} messageResponse
// @Router      /users/archive [patch]
// @Security    Bearer
func (r *userRoutes) archive(c *gin.Context) {
	var req entity.ArchiveUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.validator.TranslateAll(err))
		return
	}

	id, err := getUserID(c)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err)
		return
	}

	if err = r.userService.Archive(c.Request.Context(), id, *req.IsArchive); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary     Partial update
// @Description Update user data partially
// @ID         	update
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body entity.PartialUpdateRequest true "Provide at least one user field to update user data"
// @Success     204
// @Failure     401 {object} messageResponse
// @Failure     400 {object} messageResponse
// @Failure     500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [patch]
// @Security    Bearer
func (r *userRoutes) partialUpdate(c *gin.Context) {
	var req entity.PartialUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.validator.TranslateAll(err))
		return
	}

	id, err := getUserID(c)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err)
		return
	}

	if err = r.userService.PartialUpdate(c.Request.Context(), id, req); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary     Get user data
// @Description Receive user data
// @ID          get
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Failure     401 {object} messageResponse
// @Success     200 {object} entity.User
// @Failure     400 {object} messageResponse
// @Failure     500 {object} messageResponse
// @Router      /users [get]
// @Security    Bearer
func (r *userRoutes) getUser(c *gin.Context) {
	id, err := getUserID(c)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err)
		return
	}

	u, err := r.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, u)
}
