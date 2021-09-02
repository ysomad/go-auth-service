package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type userRoutes struct {
	userService service.User
}

// Requests and responses
type (
	createUserRequest struct {
		Email    string `json:"email"    example:"user@mail.com" binding:"required"`
		Password string `json:"password" example:"secret"        binding:"required"`
	}

	createUserResponse struct {
		ID        int       `json:"id"`
		Email     string    `json:"email"      example:"user@mail.com"`
		CreatedAt time.Time `json:"created_at" example:"2021-08-31T16:55:18.080768Z"`
	}

	archiveUserRequest struct {
		Password string `json:"password" example:"secret" binding:"required"`
	}
)

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
// @Param       request body createUserRequest true "Register a new user"
// @Param       name path string true "User id"
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

	// Validate User struct
	// TODO: refactor struct validation
	validate := validator.New()
	if err := validate.Struct(user); err != nil {

		// Init translator
		eng := en.New()
		uni := ut.New(eng, eng)
		trans, _ := uni.GetTranslator("en")

		// Register translations
		if regErr := enTranslations.RegisterDefaultTranslations(validate, trans);
		regErr != nil {
			abortWithError(c, http.StatusBadRequest, err)
			return
		}

		// Get validation errors and translate em
		validationErrs := err.(validator.ValidationErrors)
		translatedErrs := validationErrs.Translate(trans)

		// Format translated validation errors
		formattedErrs := make(validator.ValidationErrorsTranslations, len(translatedErrs))
		for k, v := range translatedErrs {
			k = strings.Split(k, ".")[1]
			words := strings.Fields(v)[1:]
			formattedErrs[strings.ToLower(k)] = strings.Join(words, " ")
		}

		abortWithValidationError(c, http.StatusUnprocessableEntity, err, formattedErrs)
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
// @Param       request body archiveUserRequest true "Archive user"
// @Param		id path int required "User ID"
// @Success     200
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /users/{id} [delete].
func (r *userRoutes) archive(c *gin.Context) {
	var request archiveUserRequest

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
		ID: id,
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
