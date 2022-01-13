package validation

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Gin interface {
	TranslateError(err error) map[string]string
}

type ginValidator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewGinValidator() (*ginValidator, error) {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, errors.New("validation translator not found")
	}

	return &ginValidator{
		validate: binding.Validator.Engine().(*validator.Validate),
		trans:    trans,
	}, nil
}

func (v *ginValidator) register() error {
	if err := enTranslations.RegisterDefaultTranslations(v.validate, v.trans); err != nil {
		return err
	}

	return nil
}

// TranslateError returns translated validation errors received from gin.c.ShouldBindJSON err
func (v *ginValidator) TranslateError(err error) map[string]string {
	_ = v.register()

	errs := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		errs[field] = err.Translate(v.trans) + "."
	}

	return errs
	// return err.(validator.ValidationErrors).Translate(v.trans)
}