package ports_authentications

import (
	"OIDCAuthenticator/internal/core/dto"
	"context"
	"crypto/rsa"
)

type RsaKeyStoreService interface {
	GetActiveKey(ctx context.Context) (string, *rsa.PrivateKey, error)
	GetPublicKeys(ctx context.Context) ([]dto.RsaKeyResult, error)
}
