package store

import "github.com/ysomad/go-auth-service/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	Get(int) (*model.User, error)
	GetByEmail(string) (*model.User, error)
}
