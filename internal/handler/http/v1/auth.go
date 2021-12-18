package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/domain"
	"github.com/ysomad/go-auth-service/internal/service"
	"github.com/ysomad/go-auth-service/pkg/validation"
)

type authHandler struct {
	validation.Validator
	session service.Session
}

func newAuthHandler(handler *gin.RouterGroup, v validation.Validator, s service.Session) {
	h := &authHandler{v, s}

	g := handler.Group("/auth")
	{
		g.POST("login", h.login)

		authenticated := g.Group("/", sessionMiddleware(s))
		{
			authenticated.POST("logout", h.logout)
			authenticated.POST("token", h.token)
		}
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=128"`
}

func (h *authHandler) login(c *gin.Context) {
	// TODO: fix
	var r loginRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithValidationError(c, http.StatusUnprocessableEntity, h.TranslateError(err))
		return
	}

	d, err := domain.NewDevice(c.Request.Header.Get("User-Agent"), c.ClientIP())
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	s, err := h.session.EmailLogin(c.Request.Context(), r.Email, r.Password, d)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	// Set httponly secure cookie with session id
	// TODO: refactor
	c.SetCookie("id", s.ID, s.TTL, "v1", "", true, true)

	c.Status(http.StatusOK)
}

func (h *authHandler) logout(c *gin.Context) {
	panic("implement")

	c.Status(http.StatusNoContent)
}

func (h *authHandler) token(c *gin.Context) {
	panic("implement")

	var t string

	c.JSON(http.StatusOK, t)
}
