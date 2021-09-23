package v1

import (
	"github.com/ysomad/go-auth-service/pkg/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service"
)

type userRoutes struct {
	translator  validation.Translator
	userService service.User
}

func newUserRoutes(handler *gin.RouterGroup, t validation.Translator, u service.User) {
	r := &userRoutes{t, u}

	h := handler.Group("/users")
	{
		h.PATCH(":id/archive", r.archive)
		h.GET(":id", r.getByID)
		h.PATCH(":id", r.partialUpdate)
		h.POST("", r.signUp)
	}
}

// @Summary     Sign up
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
		abortWithValidationErr(c, http.StatusUnprocessableEntity, r.translator.All(err))
		return
	}

	if err := r.userService.SignUp(c.Request.Context(), req); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary     Archive or restore User
// @Description Archive or restore User
// @ID          archive
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param		id path int required "User ID"
// @Param       request body entity.ArchiveUserRequest true "To archive or restore a user is_archive should be provided"
// @Success     204
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/{id}/archive [patch]
func (r *userRoutes) archive(c *gin.Context) {
	var req entity.ArchiveUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationErr(c, http.StatusUnprocessableEntity, r.translator.All(err))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
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
// @Param		id path int required "User ID"
// @Success     204
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/{id} [patch]
func (r *userRoutes) partialUpdate(c *gin.Context) {
	var req entity.PartialUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationErr(c, http.StatusUnprocessableEntity, r.translator.All(err))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	if err = r.userService.PartialUpdate(c.Request.Context(), id, req); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary     Get
// @Description Receive user data by id
// @ID          get
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param		id path int required "User ID"
// @Success     200 {object} entity.User
// @Failure     400 {object} messageResponse
// @Router      /users/{id} [get]
func (r *userRoutes) getByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	u, err := r.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, u)
}
