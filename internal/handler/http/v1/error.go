package v1

import (
	"github.com/gin-gonic/gin"
)

type messageResponse struct {
	Message string `json:"message"`
}

func abortWithError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, messageResponse{err.Error()})
}
