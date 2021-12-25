package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// TODO: refactor ERRORS or MOVE to pkg httperror

var (
	ErrLogin = errors.New("an error occurred while logging in")
)

type messageResponse struct {
	Error string `json:"error"`
}

type validationErrorResponse struct {
	Errors map[string]string `json:"error"`
}

func abortWithError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, messageResponse{err.Error()})
}

func abortWithValidationError(c *gin.Context, code int, errs map[string]string) {
	c.AbortWithStatusJSON(code, validationErrorResponse{errs})
}
