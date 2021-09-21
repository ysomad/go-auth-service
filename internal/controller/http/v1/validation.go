package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type RequestValidator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewRequestValidator() (*RequestValidator, error) {
	trans, err := getTranslator()
	if err != nil {
		return nil, err
	}

	return &RequestValidator{
		validate: binding.Validator.Engine().(*validator.Validate),
		trans:    trans,
	}, nil
}

func getTranslator() (ut.Translator, error) {
	eng := en.New()
	uni := ut.New(eng, eng)

	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, errors.New("validation translator not found")
	}

	return trans, nil
}

func (v *RequestValidator) registerTranslations() error {
	if err := enTranslations.RegisterDefaultTranslations(v.validate, v.trans); err != nil {
		return err
	}

	return nil
}

func (v *RequestValidator) translateAll(err error) map[string]string {
	return err.(validator.ValidationErrors).Translate(v.trans)
}

func ValidRequest(c *gin.Context, req interface{}) bool {
	if validationErr := c.ShouldBindJSON(req); validationErr != nil {

		v, err := NewRequestValidator()
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, err)

			return false
		}

		if regErr := v.registerTranslations(); regErr != nil {
			abortWithError(c, http.StatusInternalServerError, regErr)

			return false
		}

		abortWithValidationError(
			c,
			http.StatusUnprocessableEntity,
			v.translateAll(validationErr),
		)

		return false
	}

	return true
}
