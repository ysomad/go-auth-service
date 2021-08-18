package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "test@mail.org",
		Password: "password",
	}
}
