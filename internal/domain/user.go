package domain

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User represents user data model
type User struct {
	ID                int       `json:"id"                     example:"1"`
	Email             string    `json:"email"                  example:"user@mail.com" validate:"omitempty,email"`
	Username          string    `json:"username,omitempty"     example:"username"      validate:"omitempty,alphanum,gte=4,lte=32"`
	Password          string    `json:"password,omitempty"     example:"secret"        validate:"omitempty,gte=6,lte=128"`
	EncryptedPassword string    `json:"-"`
	FirstName         string    `json:"first_name,omitempty"   example:"Alex"          validate:"omitempty,alpha,lte=50"`
	LastName          string    `json:"last_name,omitempty"    example:"Malykh"        validate:"omitempty,alpha,lte=50"`
	CreatedAt         time.Time `json:"created_at"             example:"2021-08-31T16:55:18.080768Z"`
	IsActive          bool      `json:"is_active,omitempty"    example:"true"`
	IsSuperuser       bool      `json:"is_superuser,omitempty" example:"false"`
}

// Data transfer objects (DTO)
type (
	CreateUserRequest struct {
		Email             string `json:"email" example:"user@mail.com" binding:"required,email,lte=255"`
		Password          string `json:"password" example:"secret" binding:"required,gte=6,lte=128"`
		ConfirmPassword   string `json:"confirm_password" example:"secret" binding:"eqfield=Password"`
		EncryptedPassword string `json:"-"`
	}

	CreateUserResponse struct {
		ID        int       `json:"id"         example:"1"`
		Email     string    `json:"email"      example:"user@mail.com"`
		CreatedAt time.Time `json:"created_at" example:"2021-08-31T16:55:18.080768Z"`
	}

	ArchiveUserRequest struct {
		ID        int   `json:"-" example:"1" binding:"numeric,omitempty"`
		IsArchive *bool `json:"is_archive" example:"false" binding:"required"`
	}

	ArchiveUserResponse struct {
		ID        int       `json:"id" example:"1"`
		IsArchive bool      `json:"is_archive" example:"false"`
		UpdatedAt time.Time `json:"updated_at" example:"2021-08-31T16:55:18.080768Z"`
	}

	UpdateUserRequest struct {
		Username  string `json:"username"   example:"username" binding:"required"`
		FirstName string `json:"first_name" example:"Alex"     binding:"required"`
		LastName  string `json:"last_name"  example:"Malykh"   binding:"required"`
	}
)

func (u *CreateUserRequest) Sanitize() {
	u.Password = ""
	u.SetEncryptedPassword("")
}

func (u *CreateUserRequest) SetEncryptedPassword(p string) {
	u.EncryptedPassword = p
}

// EncryptPassword encrypts user password and write it to EncryptedPassword field of User struct
func (u *CreateUserRequest) EncryptPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.SetEncryptedPassword(string(bytes))
	return err
}

// CompareHashAndPassword compares received password from client with hashed password in db
func (u *CreateUserRequest) CompareHashAndPassword(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
	return err == nil
}
