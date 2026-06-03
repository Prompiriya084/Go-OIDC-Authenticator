package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type AuthSessionRepository interface {
	GenericRepository[domain_entities.AuthSession, domain_entities.AuthSessionFilter]
}
