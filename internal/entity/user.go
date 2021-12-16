package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ysomad/go-auth-service/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserUniqueViolation   = errors.New("user with given email already exist")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserEmptyCredentials  = errors.New("empty email, username or password")
	ErrUserIncorrectPassword = errors.New("incorrect password")
)

// User represents user data model
type User struct {
	ID          string    `json:"id" example:"c84f18a2-c6c7-4850-be15-93f9cbaef3b3"`
	Email       string    `json:"email" example:"user@mail.com"`
	Username    string    `json:"username" example:"username"`
	Password    string    `json:"-" example:"secret"`
	CreatedAt   time.Time `json:"createdAt" example:"2021-08-31T16:55:18.080768Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2021-08-31T16:55:18.080768Z"`
	IsActive    bool      `json:"isActive" example:"true"`
	IsArchive   bool      `json:"isArchive" example:"false"`
	IsSuperuser bool      `json:"-"`
}

func (u User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", ErrUserIncorrectPassword)
	}

	return nil
}

// UserSensitiveData is a data transfer object contains sensitive data.
type UserSensitiveData struct {
	email        string
	username     string
	passwordHash string
}

func (u *UserSensitiveData) Email() string {
	return u.email
}

func (u *UserSensitiveData) Username() string {
	return u.username
}

func (u *UserSensitiveData) PasswordHash() string {
	return u.passwordHash
}

func NewUserSensitiveData(email, hash string) (UserSensitiveData, error) {
	if email == "" || hash == "" {
		return UserSensitiveData{}, ErrUserEmptyCredentials
	}

	emailUsername := strings.Split(email, "@")[0]

	return UserSensitiveData{
		email:        email,
		username:     fmt.Sprintf("%s_%s", emailUsername, util.RandomString(4)),
		passwordHash: hash,
	}, nil
}
