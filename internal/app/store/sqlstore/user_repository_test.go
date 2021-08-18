package sqlstore_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/ysomad/go-auth-service/internal/app/model"
	"github.com/ysomad/go-auth-service/internal/app/store"
	"github.com/ysomad/go-auth-service/internal/app/store/sqlstore"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_Get(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	testUser := model.TestUser(t)
	s.User().Create(testUser)
	u, err := s.User().Get(testUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

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
