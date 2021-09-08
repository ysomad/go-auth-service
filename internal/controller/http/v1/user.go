package v1

import (
	"github.com/ysomad/go-auth-service/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service"
)

type userRoutes struct {
	l logger.Interface
	u service.User
}

func newUserRoutes(handler *gin.RouterGroup, l logger.Interface, u service.User) {
	r := &userRoutes{l, u}

	h := handler.Group("/users")
	{
		h.PATCH(":id/archive", r.archive)
		h.GET(":id", r.getByID)
		h.PATCH(":id", r.update)
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
// @Success     200 {object} entity.User
// @Failure     400 {object} messageResponse
// @Failure     500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [post]
func (r *userRoutes) signUp(c *gin.Context) {
	var req entity.CreateUserRequest

	if !ValidDTO(c, &req) {
		return
	}

	resp, err := r.u.SignUp(c.Request.Context(), &req)
	if err != nil {
		r.l.Error(err, "http - v1 - signUp")
		abortWithError(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, resp)
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

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - archive")
		abortWithError(c, http.StatusBadRequest, err)

		return
	}

	req.ID = id

	if !ValidDTO(c, &req) {
		return
	}

	if err = r.u.Archive(c.Request.Context(), &req); err != nil {
		r.l.Error(err, "http - v1 - archive")
		abortWithError(c, http.StatusBadRequest, err)

		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// @Summary     Partial update
// @Description Update user data partially
// @ID         	update
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body entity.UpdateUserRequest true "Provide at least one user field to update user data"
// @Failure		422 {object} validationErrorResponse
// @Param		id path int required "User ID"
// @Success     200 {object} entity.User
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/{id} [patch]
func (r *userRoutes) update(c *gin.Context) {
	var req entity.UpdateUserRequest

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - update")
		abortWithError(c, http.StatusBadRequest, err)

		return
	}

	req.ID = id

	if !ValidDTO(c, &req) {
		return
	}

	user, err := r.u.Update(c.Request.Context(), &req)
	if err != nil {
		r.l.Error(err, "http - v1 - update")
		abortWithError(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary     Get by id
// @Description Receive user data
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
		r.l.Error(err, "http - v1 - getByID")
		abortWithError(c, http.StatusBadRequest, err)

		return
	}

	user, err := r.u.GetByID(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getByID")
		abortWithError(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, user)
}
