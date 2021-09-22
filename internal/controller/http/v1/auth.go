package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service"
	"net/http"
)

type authRoutes struct {
	authService service.Auth
}

func newAuthRoutes(handler *gin.RouterGroup, a service.Auth) {
	r := &authRoutes{a}

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

	if !ValidRequest(c, &req) {
		return
	}

	// Get user data
	fingerprint, err := uuid.Parse(req.Fingerprint)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	s := entity.NewRefreshSession(
		c.Request.Header.Get("User-Agent"),
		c.ClientIP(),
		fingerprint,
	)

	// Login user
	resp, err := r.authService.Login(c.Request.Context(), req, s)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
