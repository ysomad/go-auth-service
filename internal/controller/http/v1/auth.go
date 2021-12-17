package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ysomad/go-auth-service/internal/entity"
	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

const (
	sessionIDKey = "id"
)

type authRoutes struct {
	validation.Validator
	sessionService service.Session
}

func newAuthRoutes(handler *gin.RouterGroup, v validation.Validator, a service.Session) {
	r := &authRoutes{v, a}

	h := handler.Group("/auth")
	{
		h.POST("login", r.login)
	}
}

type loginRequest struct {
	Email    string `json:"email" example:"user@mail.com" binding:"required,email,lte=255"`
	Password string `json:"password" example:"secret" binding:"required,gte=6,lte=128"`
}

// @Summary     Login
// @Description Logs in and returns authentication cookie
// @ID          authLogin
// @Tags  	    auth
// @Accept      json
// @Produce     json
// @Param       request body loginRequest true "To login user email and password should be provided."
// @Success     200
// @Header      200 {string} Set-Cookie "`id`=22KWxEi4XlPGqFrMadBFW1qEFWv; Path=v1; `HttpOnly`; `Secure`"
// @Failure     400 {object} messageResponse
// @Failure		422 {object} validationErrorResponse
// @Router      /auth/login [post]
func (r *authRoutes) login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, r.TranslateError(err))
		return
	}

	d, err := entity.NewDevice(c.Request.Header.Get("User-Agent"), c.ClientIP())
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	sess, err := r.sessionService.LoginWithEmail(c.Request.Context(), req.Email, req.Password, d)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	// Set httponly secure cookie with session id
	c.SetCookie(sessionIDKey, sess.ID, sess.TTL, "v1", "", true, true)

	c.Status(http.StatusOK)
}
