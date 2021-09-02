package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"
	"net/http"
	"time"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(handler *gin.RouterGroup, us service.User) {
	r := &userRoutes{us}

	h := handler.Group("/users")
	{
		h.POST("", r.create)
	}
}

type createUserRequest struct {
	Email    string `json:"email"    example:"user@mail.com" binding:"required"`
	Password string `json:"password" example:"secret"        binding:"required"`
}

type createUserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"      example:"user@mail.com"`
	CreatedAt time.Time `json:"created_at" example:"2021-08-31T16:55:18.080768Z"`
}

// @Summary     Create
// @Description Register a new user
// @ID          create
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body createUserRequest true "Register a new user"
// @Success     200 {object} createUserResponse
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [post].
func (r *userRoutes) create(c *gin.Context) {
	var request createUserRequest

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

	err, translatedErrs := r.userService.Create(c.Request.Context(), &user)

	// Check translated validation errors
	if translatedErrs != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, err, translatedErrs)
		return
	}

	// Check other errors
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
