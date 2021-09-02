package domain

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User represents user data model
type User struct {
	ID                int       `json:"id"                     example:"1"`
	Email             string    `json:"email"                  example:"user@mail.com" validate:"required,email"`
	Password          string    `json:"password,omitempty"     example:"secret"        validate:"required,gte=6,lte=128"`
	EncryptedPassword string    `json:"-"`
	FirstName         string    `json:"first_name,omitempty"   example:"Alex"`
	LastName          string    `json:"last_name,omitempty"    example:"Malykh"`
	CreatedAt         time.Time `json:"created_at"             example:"2021-08-31T16:55:18.080768Z"`
	IsActive          bool      `json:"is_active,omitempty"    example:"true"`
	IsSuperuser       bool      `json:"is_superuser,omitempty" example:"false"`
}

func (u *User) Sanitize() {
	u.Password = ""
}

// EncryptPassword encrypts user password and write it to EncryptedPassword field of User struct
func (u *User) EncryptPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.EncryptedPassword = string(bytes)
	return err
}

// CompareHashAndPassword compares received password from client with hashed password in db
func (u *User) CompareHashAndPassword() bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(u.Password))
	return err == nil
}
