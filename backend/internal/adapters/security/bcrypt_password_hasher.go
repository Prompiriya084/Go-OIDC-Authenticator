package adapters_security

import (
	ports_security "OIDCAuthenticator/internal/core/ports/security"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordHasher struct {
}

func NewBryptPasswordHasher() ports_security.PasswordHasher {
	return &bcryptPasswordHasher{}
}

func (h *bcryptPasswordHasher) Hash(password string) string {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedByte)
}
func (h *bcryptPasswordHasher) Verify(password string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		fmt.Println("Password validation failed!")
	} else {
		fmt.Println("Password validated successfully!")
	}
	return false
}
