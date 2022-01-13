package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/internal/service"
)

func setSessionCookie(ctx *gin.Context, c service.SessionCookie) {
	ctx.SetCookie(c.Key,c.ID,c.TTL,apiPath,c.Domain,c.Secure,c.HTTPOnly)
}