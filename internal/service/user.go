package service

import (
	"context"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"strings"

	"github.com/ysomad/go-auth-service/internal/domain"
)

type UserService struct {
	repo     UserRepo
	validate *validator.Validate
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{r, validator.New()}
}

func (s *UserService) Create(ctx context.Context, u *domain.User) (error, map[string]string) {
	// Validate User struct
	if err := s.validate.Struct(u); err != nil {

		// Init translator
		eng := en.New()
		uni := ut.New(eng, eng)
		trans, _ := uni.GetTranslator("en")

		// Register translations
		if regErr := enTranslations.RegisterDefaultTranslations(s.validate, trans); regErr != nil {
			return regErr, nil
		}

		// Get validation errors and translate em
		validationErrs := err.(validator.ValidationErrors)
		translatedErrs := validationErrs.Translate(trans)

		// Format translated validation errors
		formattedErrs := make(validator.ValidationErrorsTranslations, len(translatedErrs))
		for k, v := range translatedErrs {
			k = strings.Split(k, ".")[1]
			words := strings.Fields(v)[1:]
			formattedErrs[strings.ToLower(k)] = strings.Join(words, " ")
		}

		return err, formattedErrs
	}

	// Encrypt password
	if err := u.EncryptPassword(); err != nil {
		return err, nil
	}

	// Create a new user in database
	if err := s.repo.Create(ctx, u); err != nil {
		return err, nil
	}

	u.Sanitize()

	return nil, nil
}
