package entity

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User represents user data model
type User struct {
	ID                int       `json:"id" example:"1" binding:"numeric"`
	Email             string    `json:"email" example:"user@mail.com" binding:"email"`
	Username          *string   `json:"username,omitempty" example:"username" binding:"alphanum,gte=4,lte=32"`
	Password          string    `json:"-" example:"secret" binding:"gte=6,lte=128"`
	FirstName         *string   `json:"first_name,omitempty" example:"Alex" binding:"alpha,lte=50"`
	LastName          *string   `json:"last_name,omitempty" example:"Malykh" binding:"alpha,lte=50"`
	CreatedAt         time.Time `json:"created_at" example:"2021-08-31T16:55:18.080768Z"`
	UpdatedAt         time.Time `json:"updated_at" example:"2021-08-31T16:55:18.080768Z"`
	IsActive          bool      `json:"is_active" example:"true"`
	IsArchive         bool      `json:"is_archive" example:"false"`
	IsSuperuser       bool      `json:"-"`
}

// Data transfer objects (DTO)
type (
	CreateUserRequest struct {
		Email           string `json:"email" example:"user@mail.com" binding:"required,email,lte=255"`
		Password        string `json:"password" example:"secret" binding:"required,gte=6,lte=128"`
		ConfirmPassword string `json:"confirm_password" example:"secret" binding:"required,eqfield=Password"`
	}

	ArchiveUserRequest struct {
		IsArchive *bool `json:"is_archive" example:"false" binding:"required"`
	}

	UpdateUserRequest struct {
		ID        int    `json:"-" example:"1" binding:"required,numeric"`
		Username  string `json:"username" example:"username" binding:"omitempty,alphanum,gte=4,lte=32"`
		FirstName string `json:"first_name" example:"Alex"  binding:"omitempty,alpha,lte=50"`
		LastName  string `json:"last_name" example:"Malykh" binding:"omitempty,alpha,lte=50"`
	}
)

// EncryptPassword ...
func EncryptPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

// CompareHashAndPassword compares received password from client with hashed password in db
func (u *User) CompareHashAndPassword(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
	return err == nil
}
