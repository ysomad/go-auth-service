package domain

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserUniqueViolation   = errors.New("user with given email already exist")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserEmptyCredentials  = errors.New("empty email, username or password")
	ErrUserIncorrectPassword = errors.New("incorrect password")
)

// Account represents user data model
type Account struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	IsArchive    bool      `json:"isArchive"`
}

func (a Account) CompareHashAndPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)); err != nil {
		// TODO: return generic err pkg/httperror
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", ErrUserIncorrectPassword)
	}

	return nil
}

// AccountCredentials is a data transfer object contains sensitive data.
type AccountCredentials struct {
	email        string
	passwordHash string
}

func (a *AccountCredentials) Email() string {
	return a.email
}

func (a *AccountCredentials) PasswordHash() string {
	return a.passwordHash
}

func NewAccountCredentials(email, hash string) (AccountCredentials, error) {
	if email == "" || hash == "" {
		// TODO: return generic err pkg/httperror
		return AccountCredentials{}, ErrUserEmptyCredentials
	}

	return AccountCredentials{
		email:        email,
		passwordHash: hash,
	}, nil
}
