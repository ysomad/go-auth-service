package entity

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// User error messages
var (
	UserUniqueViolationErr = errors.New("user with given email or username already exists")
	UserNotFoundErr        = errors.New("user not found")
	UserIncorrectErr       = errors.New("incorrect email or password")
	PartialUpdateErr       = errors.New("provide at least one field to update resource partially")
)

// User represents user data model
type User struct {
	ID          uuid.UUID `json:"id" example:"c84f18a2-c6c7-4850-be15-93f9cbaef3b3" binding:"uuid4"`
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

// CreateUserRequest represents request DTO for user sign up
type CreateUserRequest struct {
	Email           string `json:"email" example:"user@mail.com" binding:"required,email,lte=255"`
	Password        string `json:"password" example:"secret" binding:"required,gte=6,lte=128"`
	ConfirmPassword string `json:"confirmPassword" example:"secret" binding:"required,eqfield=Password"`
}

// ArchiveUserRequest represents request DTO for archive or restore user operation
type ArchiveUserRequest struct {
	IsArchive *bool `json:"isArchive" example:"false" binding:"required"`
}

// PartialUpdateRequest represents request DTO for user partial update
type PartialUpdateRequest struct {
	Username  string `json:"username" example:"username" binding:"omitempty,alphanum,gte=4,lte=32"`
	FirstName string `json:"firstName" example:"Alex"  binding:"omitempty,alpha,lte=50"`
	LastName  string `json:"lastName" example:"Malykh" binding:"omitempty,alpha,lte=50"`
}
