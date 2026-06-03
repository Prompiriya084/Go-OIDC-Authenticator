package ports_repositories

import domain_entities "OIDCAuthenticator/internal/core/domain/entities"

type AuthCodeRepository interface {
	GenericRepository[domain_entities.AuthCode, domain_entities.AuthCodeFilter]
}
