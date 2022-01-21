package domain

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ysomad/go-auth-service/pkg/apperrors"
	"github.com/ysomad/go-auth-service/pkg/utils"
)

// Account represents user data model
type Account struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Archive      bool      `json:"archive"`
	Verified     bool      `json:"verified"`
}

func (a *Account) GeneratePasswordHash() error {
	b, err := bcrypt.GenerateFromPassword([]byte(a.Password), 11)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", apperrors.ErrAccountPasswordNotGenerated)
	}

	a.PasswordHash = string(b)

	return nil
}

func (a *Account) CompareHashAndPassword() error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(a.Password)); err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", apperrors.ErrAccountIncorrectPassword)
	}

	return nil
}

func (a *Account) RandomPassword() {
	a.Password = utils.RandomSpecialString(16)
}
