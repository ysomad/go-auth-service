package sqlstore

import (
	"database/sql"
	"github.com/ysomad/go-auth-service/internal/app/model"
	"github.com/ysomad/go-auth-service/internal/app/store"
)

type UserRepository struct {
	store *Store
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

	// Insert user to db
	return r.store.db.QueryRow(
		"insert into users (email, password) values ($1, $2) returning id",
		u.Email, u.EncryptedPassword,
	).Scan(&u.ID)
}

func (r *UserRepository) Get(id int) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRow(
		"select id, email, password from users where id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// GetByEmail searches for user with particular email address
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRow(
		"select id, email, password from users where email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
