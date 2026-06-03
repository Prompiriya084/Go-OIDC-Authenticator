package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type ClientScopeRepository interface {
	GenericRepository[domain_entities.ClientScope, domain_entities.ClientScopeFilter]
}
