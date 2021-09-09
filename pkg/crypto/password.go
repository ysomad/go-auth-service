package crypto

import "golang.org/x/crypto/bcrypt"

type PasswordService struct{}

func NewPasswordService() PasswordService {
	return PasswordService{}
}

// Encrypt returns hashed password string
func (p PasswordService) Encrypt(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// Compare compares hashed password with string password
func (p PasswordService) Compare(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
