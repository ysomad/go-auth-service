package entity

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User represents user data model
type User struct {
	ID          int       `json:"id" example:"1" binding:"numeric"`
	Email       string    `json:"email" example:"user@mail.com" binding:"email"`
	Username    *string   `json:"username,omitempty" example:"username" binding:"alphanum,gte=4,lte=32"`
	Password    string    `json:"-" example:"secret" binding:"gte=6,lte=128"`
	FirstName   *string   `json:"firstName,omitempty" example:"Alex" binding:"alpha,lte=50"`
	LastName    *string   `json:"lastName,omitempty" example:"Malykh" binding:"alpha,lte=50"`
	CreatedAt   time.Time `json:"createdAt" example:"2021-08-31T16:55:18.080768Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2021-08-31T16:55:18.080768Z"`
	IsActive    bool      `json:"isActive" example:"true"`
	IsArchive   bool      `json:"isArchive" example:"false"`
	IsSuperuser bool      `json:"-"`
}

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
