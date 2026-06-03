package ports_security

type PasswordHasher interface {
	Hash(password string) string
	Verify(password string, passwordHash string) bool
}
