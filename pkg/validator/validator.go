package validator

import (
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validate *v10.Validate
	trans    ut.Translator
}

func New() *Validator {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")

	return &Validator{v10.New(), trans}
}

// ValidateStruct is a shortcut for validator.Validate.Struct method
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// Fmt formats translated validation errors by go-playground/validator package
func (v *Validator) Fmt(errs map[string]string) map[string]string {
	fmtErrs := make(map[string]string, len(errs))
	for k, v := range errs {
		k = strings.Split(k, ".")[1]
		words := strings.Fields(v)[1:]
		fmtErrs[strings.ToLower(k)] = strings.Join(words, " ")
	}

	return fmtErrs
}

// TranslateAll translates all validation errors received from go-playground/validator package
func (v *Validator) TranslateAll(err error) (map[string]string, error) {
	if regErr := enTranslations.RegisterDefaultTranslations(v.validate, v.trans); regErr != nil {
		return nil, regErr
	}

	return err.(v10.ValidationErrors).Translate(v.trans), nil
}
