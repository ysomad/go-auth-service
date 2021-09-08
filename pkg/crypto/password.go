package crypto

import "golang.org/x/crypto/bcrypt"

// EncryptPassword returns hashed password string
func EncryptPassword(p string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), cost)
	return string(bytes), err
}

// CompareHashAndPassword compares hashed password with string password
func CompareHashAndPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
