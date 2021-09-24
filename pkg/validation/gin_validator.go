package validation

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	ValidateVar(val interface{}, tag string) error
	TranslateAll(err error) map[string]string
}

type GinValidator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewGinValidator() (*GinValidator, error) {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, errors.New("validation translator not found")
	}

	return &GinValidator{
		validate: binding.Validator.Engine().(*validator.Validate),
		trans:    trans,
	}, nil
}

func (v *GinValidator) register() error {
	if err := enTranslations.RegisterDefaultTranslations(v.validate, v.trans); err != nil {
		return err
	}

	return nil
}

// TranslateAll returns translated validation errors received from gin.c.ShouldBindJSON err
func (v *GinValidator) TranslateAll(err error) map[string]string {
	v.register()

	return err.(validator.ValidationErrors).Translate(v.trans)
}

func (v *GinValidator) ValidateVar(val interface{}, tag string) error {
	return v.validate.Var(val, tag)
}
