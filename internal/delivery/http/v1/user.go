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
		h.PATCH(":id/state", r.updateState)
		h.PUT(":id", r.update)
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
// @Success     201 {object} domain.CreateUserResponse
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

	c.JSON(http.StatusCreated, resp)
}

// @Summary     Update state
// @Description Update user state
// @ID          state
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param		id path int required "User ID"
// @Param       request body domain.UpdateStateUserRequest true "To change user state is_archive should be provided"
// @Success     204
// @Failure     400 {object} messageResponse
// @Router      /users/{id}/state [patch].
func (r *userRoutes) updateState(c *gin.Context) {
	var request domain.UpdateStateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	user := domain.User{
		ID:       id,
		IsActive: *request.IsActive,
	}

	if err = r.userService.UpdateState(c.Request.Context(), &user); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// @Summary     Update
// @Description Update user data
// @ID          update
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body domain.UpdateUserRequest true "All required fields should be provided"
// @Failure		422 {object} validationErrorResponse
// @Param		id path int required "User ID"
// @Success     204
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/{id} [put].
func (r *userRoutes) update(c *gin.Context) {
	var request domain.UpdateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	user := domain.User{
		ID:        id,
		Username:  request.Username,
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}

	if err = r.userService.Update(c.Request.Context(), &user); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
