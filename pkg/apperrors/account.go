package apperrors

import "errors"

var (
	ErrAccountAlreadyExist             = errors.New("account with given email or username already exist")
	ErrAccountNotFound                 = errors.New("account not found")
	ErrAccountNotArchived              = errors.New("account cannot be archived")
	ErrAccountIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrAccountPasswordNotGenerated     = errors.New("password generation error")
	ErrAccountIncorrectPassword        = errors.New("incorrect password")
	ErrAccountContextNotFound          = errors.New("account not found in context")
)
