package entity

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrUserUniqueViolation    = errors.New("user with given email or username already exists")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserInvalidCredentials = errors.New("invalid email or password")
	ErrPartialUpdate          = errors.New("provide at least one field to update resource partially")
)

// User represents user data model
type User struct {
	ID          uuid.UUID `json:"id" example:"c84f18a2-c6c7-4850-be15-93f9cbaef3b3"`
	Email       string    `json:"email" example:"user@mail.com"`
	Username    *string   `json:"username,omitempty" example:"username"`
	Password    string    `json:"-" example:"secret"`
	FirstName   *string   `json:"firstName,omitempty" example:"Alex"`
	LastName    *string   `json:"lastName,omitempty" example:"Malykh"`
	CreatedAt   time.Time `json:"createdAt" example:"2021-08-31T16:55:18.080768Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2021-08-31T16:55:18.080768Z"`
	IsActive    bool      `json:"isActive" example:"true"`
	IsArchive   bool      `json:"isArchive" example:"false"`
	IsSuperuser bool      `json:"-"`
}

func (u *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return ErrUserInvalidCredentials
	}

	return nil
}

type UserCredentialsDTO struct {
	Email    string
	Password string
}

type UserPartialUpdateDTO struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
}

type UpdateColumns map[string]interface{}

func (c UpdateColumns) Validate() error {
	for k, v := range c {
		if v == "" || v == nil {
			delete(c, k)
		}
	}

	if len(c) == 0 {
		return ErrPartialUpdate
	}

	return nil
}
