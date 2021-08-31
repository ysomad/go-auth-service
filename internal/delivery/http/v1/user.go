package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"
	"net/http"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(handler *gin.RouterGroup, us service.User) {
	r := &userRoutes{us}

	h := handler.Group("/users")
	{
		h.POST("/create", r.create)
	}
}

type createUserRequest struct {
	Email    string `json:"email"    binding:"required" example:"user@mail.com"`
	Password string `json:"password" binding:"required" example:"secret"`
}

// @Summary     Create
// @Description Register a new user
// @ID          create
// @Tags  	    create
// @Accept      json
// @Produce     json
// @Param       request body createUserRequest true "Register a new user"
// @Success     200 {object} domain.User
// @Failure     400 {object} response
// @Router      /users/create [post].
func (r *userRoutes) create(c *gin.Context) {
	var request createUserRequest

	// Validate request body
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, err, "invalid request body")
		return
	}

	// Create user
	user := domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	if err := r.userService.Create(c.Request.Context(), user); err != nil {
		errorResponse(c, http.StatusBadRequest, err, "user service error")
		return
	}

	c.JSON(http.StatusOK, user)
}
