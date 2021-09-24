package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
	"net/http"
)

type authRoutes struct {
	validator   validation.Validator
	authService service.Auth
}

func newAuthRoutes(handler *gin.RouterGroup, t validation.Validator, a service.Auth) {
	r := &authRoutes{t, a}

	h := handler.Group("/auth")
	{
		h.POST("login", r.login)
	}
}

// @Summary     Login
// @Description Create access and refresh tokens using user email and password
// @ID          login
// @Tags  	    Authorization
// @Accept      json
// @Produce     json
// @Param       request body entity.LoginRequest true "To login user email, password and fingerprint as uuid v4 type should be provided"
// @Success     200 {object} entity.LoginResponse
// @Failure     400 {object} messageResponse
// @Failure     500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /auth/login [post]
func (r *authRoutes) login(c *gin.Context) {
	var req entity.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.validator.TranslateAll(err))
		return
	}

	// Get user data
	fingerprint, err := uuid.Parse(req.Fingerprint)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	// Login user
	resp, err := r.authService.Login(c.Request.Context(), req, entity.Session{
		UserAgent:   c.Request.Header.Get("User-Agent"),
		UserIP:      c.ClientIP(),
		Fingerprint: fingerprint,
	})
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	// Set httponly secure cookie with refresh token to reuse it in web applications
	c.SetCookie(
		"refreshToken",
		resp.RefreshToken.String(),
		resp.ExpiresIn,
		"/v1/auth",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, resp)
}
