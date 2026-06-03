package adapters_authentications

import (
	ports_authentications "OIDCAuthenticator/internal/core/ports/authentications"
	"crypto/rand"
	"encoding/base64"
)

type randomNumberGeneratorImpl struct{}

func NewRandomNumberGenerator() ports_authentications.RandomNumberGenerator {
	return &randomNumberGeneratorImpl{}
}

func (r *randomNumberGeneratorImpl) ToBase64String() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
