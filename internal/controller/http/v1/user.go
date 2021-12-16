package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type userRoutes struct {
	validation.Validator
	userService service.User
}

func newUserRoutes(handler *gin.RouterGroup, v validation.Validator, u service.User, s service.Session) {
	r := &userRoutes{v, u}

	h := handler.Group("/users")
	{
		h.POST("", r.register)

		authenticated := h.Group("/", sessionMiddleware(s))
		{
			authenticated.GET(":userID", r.getUser)
			authenticated.PATCH("archive/:userID", r.archive)
		}
	}
}

type userCreateRequest struct {
	Email           string `json:"email" example:"user@mail.com" binding:"required,email,lte=255"`
	Password        string `json:"password" example:"secret" binding:"required,gte=6,lte=128"`
	ConfirmPassword string `json:"confirmPassword" example:"secret" binding:"required,eqfield=Password"`
}

// @Summary     Register
// @Description Register a new user with email and password
// @ID          userCreate
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body userCreateRequest true "To register a new user email and password should be provided"
// @Success     200 {object} entity.User
// @Failure     400,500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users [post]
func (r *userRoutes) register(c *gin.Context) {
	var req userCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.TranslateError(err))
		return
	}

	user, err := r.userService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
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
// @Param       user_id path string true "User ID"
// @Router      /users/{user_id}/archive [patch]
func (r *userRoutes) archive(c *gin.Context) {
	var req userArchiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.TranslateError(err))
		return
	}

	userID, found := c.Params.Get("userID")
	if !found {
		abortWithError(c, http.StatusUnauthorized, errors.New("URL param is empty"))
		return
	}

	if err := r.userService.Archive(c.Request.Context(), userID, *req.IsArchive); err != nil {
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
// @Param       user_id path string true "User ID"
// @Router      /users/{user_id} [get]
func (r *userRoutes) getUser(c *gin.Context) {
	userID, found := c.Params.Get("userID")
	if !found {
		abortWithError(c, http.StatusUnauthorized, entity.ErrUnauthorized)
		return
	}

	u, err := r.userService.FindByID(c.Request.Context(), userID)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, u)
}
