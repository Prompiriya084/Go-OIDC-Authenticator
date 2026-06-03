package ports_authentications

import (
	"context"
	"crypto/rsa"
)

type KeyResult struct {
	Kid string
	Key *rsa.PublicKey
}

type RsaKeyStoreService interface {
	GetActiveKey(ctx context.Context) (string, *rsa.PrivateKey, error)
	GetPublicKeys(ctx context.Context) ([]KeyResult, error)
}
