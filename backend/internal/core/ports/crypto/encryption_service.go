package ports_crypto

type EncryptionService interface {
	Encrypt(text string) (string, error)
	Decrypt(cipherText string) (string, error)
}
