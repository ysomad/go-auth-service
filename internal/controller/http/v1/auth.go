package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
	"net/http"
)

const (
	refreshTokenKey = "token"
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
		h.POST("refresh", r.refreshJWT)
	}
}

// @Summary     Login
// @Description Create access and refresh tokens using user email and password
// @ID          login
// @Tags  	    Auth
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
	fp, err := uuid.Parse(req.Fingerprint)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	// Login user
	resp, err := r.authService.Login(c.Request.Context(), req, entity.SessionSecurityDTO{
		UserAgent:   c.Request.Header.Get("User-Agent"),
		UserIP:      c.ClientIP(),
		Fingerprint: fp,
	})
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	// Set httponly secure cookie with refresh token to reuse it in web applications
	c.SetCookie(
		refreshTokenKey,
		resp.RefreshToken.String(),
		resp.ExpiresIn,
		"/v1/auth",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, resp)
}

type refreshJWTRequest struct {
	RefreshToken string `json:"refreshToken" example:"c84f18a2-c6c7-4850-be15-93f9cbaef3b3" binding:"required,uuid4"`
	Fingerprint  string `json:"fingerprint" example:"c84f18a2-c6c7-4850-be15-93f9cbaef3b3" binding:"required,uuid4"`
}

// @Summary     Refresh access token
// @Description Creates new access token
// @ID          refresh
// @Tags  	    Auth
// @Accept      json
// @Produce     json
// @Param       request body refreshJWTRequest true "To get new access token fingerprint and refresh token should be provided"
// @Success     200 {object} entity.LoginResponse
// @Failure     400 {object} messageResponse
// @Failure     500 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /auth/refresh [post]
func (r *authRoutes) refreshJWT(c *gin.Context) {
	var req refreshJWTRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.validator.TranslateAll(err))
		return
	}

	fp, err := uuid.Parse(req.Fingerprint)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	rt, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	resp, err := r.authService.RefreshToken(c.Request.Context(), entity.SessionSecurityDTO{
		RefreshToken: rt,
		UserAgent:    c.Request.Header.Get("User-Agent"),
		UserIP:       c.ClientIP(),
		Fingerprint:  fp,
	})
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
