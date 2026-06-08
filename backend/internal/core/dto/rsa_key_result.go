package dto

import "crypto/rsa"

type RsaKeyResult struct {
	Kid string
	Key *rsa.PublicKey
}
