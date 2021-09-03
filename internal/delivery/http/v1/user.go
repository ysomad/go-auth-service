package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validator"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(handler *gin.RouterGroup, us service.User) {
	r := &userRoutes{us}

	h := handler.Group("/users")
	{
		h.POST("", r.create)
		h.DELETE(":id", r.archive)
	}
}

// @Summary     Create
// @Description Register a new user with email and password
// @ID          create
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body domain.CreateUserRequest true "Register a new user"
// @Success     200 {object} domain.CreateUserResponse
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [post].
func (r *userRoutes) create(c *gin.Context) {
	var request domain.CreateUserRequest

	// Validate request body
	if err := c.ShouldBindJSON(&request); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	// Pre-populate User struct
	user := domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	// Validate User struct
	v := validator.New()
	if err := v.ValidateStruct(user); err != nil {
		// Translate validation errors
		translatedErrs, translateErr := v.TranslateAll(err)
		if translateErr != nil {
			abortWithError(c, http.StatusBadRequest, translateErr)
			return
		}

		abortWithValidationError(c, http.StatusUnprocessableEntity, err, v.Fmt(translatedErrs))
		return
	}

	if err := r.userService.Create(c.Request.Context(), &user); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary     Archive
// @Description Archive user if password is correct
// @ID          archive
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body domain.ArchiveUserRequest true "Archive user"
// @Param		id path int required "User ID"
// @Success     200
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/{id} [delete].
func (r *userRoutes) archive(c *gin.Context) {
	var request domain.ArchiveUserRequest

	// Validate request body
	if err := c.ShouldBindJSON(&request); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	// Get user id from url param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	user := domain.User{
		ID:       id,
		Password: request.Password,
	}

	// Validate struct
	// TODO: Add struct validation for archive endpoint

	// Archive user
	if err = r.userService.Archive(c.Request.Context(), &user); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
