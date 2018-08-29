package security

import "golang.org/x/crypto/bcrypt"

func GetHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func ComparePasswords(hash string, password string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return res == nil
}