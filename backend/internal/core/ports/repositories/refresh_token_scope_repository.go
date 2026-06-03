package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type RefreshTokenScopeRepository interface {
	GenericRepository[domain_entities.RefreshTokenScope, domain_entities.RefreshTokenScopeFilter]
}
