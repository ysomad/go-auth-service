package domain

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User represents user data model
type User struct {
	ID                int       `json:"id"                     example:"1"`
	Email             string    `json:"email"                  example:"user@mail.com" validate:"email"`
	Username          string    `json:"username,omitempty"     example:"username"      validate:"omitempty,alpha,gte=4,lte=32"`
	Password          string    `json:"password,omitempty"     example:"secret"        validate:"gte=6,lte=128"`
	EncryptedPassword string    `json:"-"`
	FirstName         string    `json:"first_name,omitempty"   example:"Alex"          validate:"omitempty,alpha,lte=50"`
	LastName          string    `json:"last_name,omitempty"    example:"Malykh"        validate:"omitempty,alpha,lte=50"`
	CreatedAt         time.Time `json:"created_at"             example:"2021-08-31T16:55:18.080768Z"`
	IsActive          bool      `json:"is_active,omitempty"    example:"true"`
	IsSuperuser       bool      `json:"is_superuser,omitempty" example:"false"`
}

// Requests and responses
type (
	CreateUserRequest struct {
		Email    string `json:"email"    example:"user@mail.com" binding:"required"`
		Password string `json:"password" example:"secret"        binding:"required"`
	}

	CreateUserResponse struct {
		ID        int       `json:"id"`
		Email     string    `json:"email"      example:"user@mail.com"`
		CreatedAt time.Time `json:"created_at" example:"2021-08-31T16:55:18.080768Z"`
	}

	ArchiveUserRequest struct {
		IsActive *bool `json:"is_active" example:"true" binding:"required"`
	}

	UpdateUserRequest struct {
		Username  string `json:"username"   example:"username"`
		FirstName string `json:"first_name" example:"Alex"`
		LastName  string `json:"last_name"  example:"Malykh"`
	}
)

func (u *User) Sanitize() {
	u.Password = ""
	u.setEncryptedPassword("")
}

func (u *User) setEncryptedPassword(p string) {
	u.EncryptedPassword = p
}

// EncryptPassword encrypts user password and write it to EncryptedPassword field of User struct
func (u *User) EncryptPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.setEncryptedPassword(string(bytes))
	return err
}

// CompareHashAndPassword compares received password from client with hashed password in db
func (u *User) CompareHashAndPassword(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
	return err == nil
}
