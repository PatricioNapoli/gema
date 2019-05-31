package security

import (
	"gema/server/utils"

	"golang.org/x/crypto/bcrypt"
)

// Returns the bcrypt hash based on the password.
func GetHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	utils.Handle(err)
	return string(hash)
}

// Verifies the password with its hash, returns true if equal, false otherwise.
func VerifyPassword(hash string, password string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return res == nil
}
