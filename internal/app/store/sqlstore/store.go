package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/ysomad/go-auth-service/internal/app/store"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New Store creation method
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User method for using UserRepository from outside
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
