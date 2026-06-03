package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type RefreshTokenRepository interface {
	GenericRepository[domain_entities.RefreshToken, domain_entities.RefreshTokenFilter]
}
