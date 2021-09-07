package domain

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
	EncryptedPassword string    `json:"-"`
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
		ID        int   `json:"-" example:"1" binding:"numeric"`
		IsArchive *bool `json:"is_archive" example:"false" binding:"required"`
	}

	ArchiveUserResponse struct {
		ID        int       `json:"id" example:"1"`
		IsArchive bool      `json:"is_archive" example:"false"`
		UpdatedAt time.Time `json:"updated_at" example:"2021-08-31T16:55:18.080768Z"`
	}

	UpdateUserRequest struct {
		ID        int    `json:"-" example:"1" binding:"required,numeric"`
		Username  string `json:"username" example:"username" binding:"omitempty,alphanum,gte=4,lte=32"`
		FirstName string `json:"first_name" example:"Alex"  binding:"omitempty,alpha,lte=50"`
		LastName  string `json:"last_name" example:"Malykh" binding:"omitempty,alpha,lte=50"`
	}
)

func (u *User) Sanitize() {
	u.Password = ""
	u.SetEncryptedPassword("")
}

func (u *User) SetEncryptedPassword(p string) {
	u.EncryptedPassword = p
}

// EncryptPassword encrypts user password and write it to EncryptedPassword field of User struct
func (u *User) EncryptPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.SetEncryptedPassword(string(bytes))
	return err
}

// CompareHashAndPassword compares received password from client with hashed password in db
func (u *User) CompareHashAndPassword(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
	return err == nil
}
