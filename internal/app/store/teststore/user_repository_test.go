package teststore_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/ysomad/go-auth-service/internal/app/model"
	"github.com/ysomad/go-auth-service/internal/app/store"
	"github.com/ysomad/go-auth-service/internal/app/store/teststore"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_Get(t *testing.T) {
	s := teststore.New()
	testUser := model.TestUser(t)
	s.User().Create(testUser)
	u, err := s.User().Get(testUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	s := teststore.New()

	// Test with not existing user with email
	email := "test@mail.org"
	_, err := s.User().GetByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	// Test with existing user
	testUser := model.TestUser(t)
	testUser.Email = email
	s.User().Create(testUser)
	u, err := s.User().GetByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
