package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type ClientRepository interface {
	GenericRepository[domain_entities.Client, domain_entities.ClientFilter]
}
