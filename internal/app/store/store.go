package store

type Store interface {
	User() UserRepository
}
