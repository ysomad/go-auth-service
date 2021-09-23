package validation

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Translator interface {
	All(err error) map[string]string
}

type GinTranslator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewGinTranslator() (*GinTranslator, error) {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, errors.New("validation translator not found")
	}

	return &GinTranslator{
		validate: binding.Validator.Engine().(*validator.Validate),
		trans:    trans,
	}, nil
}

func (t *GinTranslator) register() error {
	if err := enTranslations.RegisterDefaultTranslations(t.validate, t.trans); err != nil {
		return err
	}

	return nil
}

// All returns translated validation errors received from gin.c.ShouldBindJSON err
func (t *GinTranslator) All(err error) map[string]string {
	t.register()

	return err.(validator.ValidationErrors).Translate(t.trans)
}
