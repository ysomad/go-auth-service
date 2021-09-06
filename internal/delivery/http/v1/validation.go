package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type DTOValidator struct {
	v *validator.Validate
	t ut.Translator
}

func New() (*DTOValidator, error) {
	trans, err := getTranslator()
	if err != nil {
		return nil, err
	}

	return &DTOValidator{
		v: binding.Validator.Engine().(*validator.Validate),
		t: trans,
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
		//k = strings.Split(k, ".")[1]
		words := strings.Fields(v)[1:]
		fmtErrs[k] = strings.Join(words, " ")
	}

	return fmtErrs
}

func (dtov *DTOValidator) registerTranslations() error {
	if err := enTranslations.RegisterDefaultTranslations(dtov.v, dtov.t); err != nil {
		return err
	}

	return nil
}

func (dtov *DTOValidator) translateAll(err error) map[string]string {
	return err.(validator.ValidationErrors).Translate(dtov.t)
}

// ValidDTO validates DTO and if it's not valid return translated validation errors to client
func ValidDTO(c *gin.Context, dto interface{}) bool {
	if validationErr := c.ShouldBindJSON(dto); validationErr != nil {

		dtov, err := New()
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, err)
			return false
		}

		if regErr := dtov.registerTranslations(); regErr != nil {
			abortWithError(c, http.StatusInternalServerError, regErr)
			return false
		}

		abortWithValidationError(
			c,
			http.StatusUnprocessableEntity,
			validationErr,
			fmtValidationErrs(dtov.translateAll(validationErr)),
		)
		return false
	}

	return true
}
