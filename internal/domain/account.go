package domain

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
)

const AccountCacheKey = "acc"

// Account represents user data model
type Account struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"passwordHash"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	IsArchive    bool      `json:"isArchive"`
}

func (a *Account) GeneratePasswordHash(password string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", apperrors.ErrAccountPasswordNotGenerated)
	}

	a.PasswordHash = string(b)

	return nil
}

func (a *Account) CompareHashAndPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)); err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", apperrors.ErrAccountIncorrectPassword)
	}

	return nil
}
