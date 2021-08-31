// Package domain defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package domain

import (
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID                int       `json:"id"         example:"1"`
	Email             string    `json:"email"      example:"user@mail.com" validate:"required,email"`
	Password          string    `json:"password"   example:"secret"        validate:"required"`
	EncryptedPassword string    `json:"-"`
	FirstName         string    `json:"first_name" example:"Alex"`
	LastName          string    `json:"last_name"  example:"Malykh"`
	CreatedAt         time.Time `json:"created_at" example:"2009-11-10 23:00:00"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *User) SetCreatedAt(t time.Time) {
	u.CreatedAt = t
}

// HashPassword hashes user password and write it to EncryptedPassword field of User struct
func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.EncryptedPassword = string(bytes)
	return err
}

func (u *User) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
