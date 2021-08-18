package teststore

import (
	"github.com/ysomad/go-auth-service/internal/app/model"
	"github.com/ysomad/go-auth-service/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	// Validate user fields before inserting to database
	if err := u.Validate(); err != nil {
		return err
	}

	// Encrypt password before inserting to db
	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

func (r *UserRepository) Get(id int) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

// GetByEmail searches for user with particular email address
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
