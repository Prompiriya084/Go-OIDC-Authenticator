package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type ClientGrantTypeRepository interface {
	GenericRepository[domain_entities.ClientGrantType, domain_entities.ClientGrantTypeFilter]
}
