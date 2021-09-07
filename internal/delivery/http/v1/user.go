package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"
	"net/http"
	"strconv"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(handler *gin.RouterGroup, us service.User) {
	r := &userRoutes{us}

	h := handler.Group("/users")
	{
		h.PATCH(":id/archive", r.archive)
		h.PATCH(":id", r.update)
		h.POST("", r.signUp)
	}
}

// @Summary     Sign Up
// @Description Register a new user with email and password
// @ID          signup
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body domain.CreateUserRequest true "To register a new user email and password should be provided"
// @Success     200 {object} domain.User
// @Failure     400 {object} messageResponse
// @Failure     500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [post].
func (r *userRoutes) signUp(c *gin.Context) {
	var req domain.CreateUserRequest

	if !ValidDTO(c, &req) {
		return
	}

	resp, err := r.userService.SignUp(c.Request.Context(), &req)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Archive user
// @Description Archive/restore user
// @ID          archive
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param		id path int required "User ID"
// @Param       request body domain.ArchiveUserRequest true "To archive/restore user is_archive boolean should be provided"
// @Success     204
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/{id}/archive [patch].
func (r *userRoutes) archive(c *gin.Context) {
	var req domain.ArchiveUserRequest

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	req.ID = id

	if !ValidDTO(c, &req) {
		return
	}

	if err = r.userService.Archive(c.Request.Context(), &req); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// @Summary     Partial Update
// @Description Update user data partially
// @ID          partialUpdate
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body domain.UpdateUserRequest true "Provide at least one user field to update user data"
// @Failure		422 {object} validationErrorResponse
// @Param		id path int required "User ID"
// @Success     200 {object} domain.User
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/{id} [patch].
func (r *userRoutes) update(c *gin.Context) {
	var req domain.UpdateUserRequest

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	req.ID = id

	if !ValidDTO(c, &req) {
		return
	}

	user, err := r.userService.Update(c.Request.Context(), &req)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
