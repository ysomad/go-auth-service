package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(handler *gin.RouterGroup, us service.User) {
	r := &userRoutes{us}

	h := handler.Group("/users")
	{
		h.PATCH(":id/archive", r.archive)
		h.PUT(":id", r.update)
		h.POST("", r.create)
	}
}

// @Summary     Create
// @Description Register a new user with email and password
// @ID          create
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body domain.CreateUserRequest true "To register a new user email and password should be provided"
// @Success     200 {object} domain.CreateUserResponse
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [post].
func (r *userRoutes) create(c *gin.Context) {
	var request domain.CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	user := domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	if !validStruct(c, user) {
		return
	}

	if err := r.userService.Create(c.Request.Context(), &user); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary     Archive
// @Description Update user state
// @ID          archive
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param		id path int required "User ID"
// @Param       request body domain.ArchiveUserRequest true "To change user state is_archive should be provided"
// @Success     200
// @Failure     400 {object} messageResponse
// @Router      /users/{id}/archive [patch].
func (r *userRoutes) archive(c *gin.Context) {
	var request domain.ArchiveUserRequest

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

	c.AbortWithStatus(http.StatusOK)
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
// @Success     200
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/{id} [put].
func (r *userRoutes) update(c *gin.Context) {
	// TODO: recreate with PATCH partial update https://play.golang.org/p/IQAHgqfBRh
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

	if !validStruct(c, user) {
		return
	}

	if err = r.userService.Update(c.Request.Context(), &user); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
