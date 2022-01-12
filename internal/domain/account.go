package domain

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
	"github.com/ysomad/go-auth-service/pkg/util"
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
	IsArchive    bool      `json:"isArchive"`
	IsVerified   bool      `json:"isVerified"`
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
	a.Password = util.RandomSpecialString(16)
}
