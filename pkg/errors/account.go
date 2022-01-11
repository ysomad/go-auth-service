package errors

import "errors"

var (
	ErrAccountAlreadyExist             = errors.New("account with given credentials already exist")
	ErrAccountNotFound                 = errors.New("account not found")
	ErrAccountNotArchived              = errors.New("account cannot be archived")
	ErrAccountIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrAccountPasswordNotGenerated     = errors.New("password generation error")
	ErrAccountIncorrectPassword        = errors.New("incorrect password")
	ErrAccountNotInContext             = errors.New("account not found in context")
)
