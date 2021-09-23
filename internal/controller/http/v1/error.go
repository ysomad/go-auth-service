package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type messageResponse struct {
	Error string `json:"error" example:"message"`
}

type validationErrResponse struct {
	Errors validator.ValidationErrorsTranslations `json:"error" example:"ModelName.FieldName:validation error message"`
}

func abortWithError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, messageResponse{err.Error()})
}

func abortWithValidationErr(c *gin.Context, code int, errs validator.ValidationErrorsTranslations) {
	c.AbortWithStatusJSON(code, validationErrResponse{errs})
}
