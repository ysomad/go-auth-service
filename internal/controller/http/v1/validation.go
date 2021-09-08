package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"net/http"
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

// fmtValidationErrs formats translated validation errors
func fmtValidationErrs(errs map[string]string) map[string]string {
	fmtErrs := make(map[string]string, len(errs))
	for k, v := range errs {
		fmtErrs[k] = v
	}

	return fmtErrs
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
			fmtValidationErrs(v.translateAll(validationErr)),
		)

		return false
	}

	return true
}

// stripNilValues removes empty strings and nil values from map https://github.com/Masterminds/squirrel/issues/66
func stripNilValues(in map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	for k, v := range in {
		if v != nil && v != "" {
			out[k] = v
		}
	}

	if len(out) == 0 {
		return nil, errors.New("provide at least one field to update resource partially")
	}

	return out, nil
}
