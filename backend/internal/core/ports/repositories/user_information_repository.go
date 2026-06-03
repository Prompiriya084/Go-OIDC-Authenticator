package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type UserInformationRepository interface {
	GenericRepository[domain_entities.UserInformation, domain_entities.UserInformationFilter]
}
