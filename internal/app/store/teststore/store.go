package teststore

import (
	"github.com/ysomad/go-auth-service/internal/app/model"
	"github.com/ysomad/go-auth-service/internal/app/store"
)

type Store struct {
	userRepository *UserRepository
}

// New Store creation method
func New() *Store {
	return &Store{}
}

// User method for using UserRepository from outside
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.userRepository
}
