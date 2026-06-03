package ports_authentications

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"context"
	"time"
)

type JwtTokenService interface {
	CreateIdToken(
		ctx context.Context,
		userId string,
		clientId string,
		nonce *string,
		userInfo *domain_entities.UserInformation,
		expiryDateUtc time.Time,
	) (string, error)
	CreateAccessToken(
		ctx context.Context,
		userId string,
		clientId string,
		audiences []string,
		scopes []string,
		expiryDateUtc time.Time,
	) (string, error)
}
