package ports_security

type Sha256Hasher interface {
	Hash(text string) string
}
