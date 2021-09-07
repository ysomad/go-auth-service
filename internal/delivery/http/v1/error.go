package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/ysomad/go-auth-service/pkg/logger"
)

type messageResponse struct {
	Error string `json:"error" example:"message"`
}

type validationErrorResponse struct {
	Errors map[string]string `json:"error" example:"ModelName.FieldName:validation error message"`
}

func logError(c *gin.Context, code int, err error, msg string) {
	logger.Error(
		err,
		msg,
		logger.Field{Key: "path", Val: c.FullPath()},
		logger.Field{Key: "request_method", Val: c.Request.Method},
		logger.Field{Key: "response_code", Val: code},
	)
}

func abortWithError(c *gin.Context, code int, err error) {
	logError(c, code, err, "http - v1 - abortWithError")
	c.AbortWithStatusJSON(code, messageResponse{err.Error()})
}

func abortWithValidationError(c *gin.Context, code int, err error, errs map[string]string) {
	logError(c, code, err, "http - v1 - abortWithValidationError")
	c.AbortWithStatusJSON(code, validationErrorResponse{errs})
}
