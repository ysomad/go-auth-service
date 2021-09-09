package crypto

//go:generate mockgen -source=crypto.go -destination=./mocks/mock.go -package=crypto_mock

type Password interface {
	Encrypt(password string, cost int) (string, error)
	Compare(hash string, password string) bool
}
