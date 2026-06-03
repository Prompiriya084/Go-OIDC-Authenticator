package ports_security

type PkceHasher interface {
	Validate(verifier string, challenge string) bool
}
