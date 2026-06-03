package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type UserMfaRepository interface {
	GenericRepository[domain_entities.UserMfa, domain_entities.UserMfaFilter]
}
