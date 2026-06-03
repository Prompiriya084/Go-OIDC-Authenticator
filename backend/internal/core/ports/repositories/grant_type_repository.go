package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type GrantTypeRepository interface {
	GenericRepository[domain_entities.GrantType, domain_entities.GrantTypeFilter]
}
