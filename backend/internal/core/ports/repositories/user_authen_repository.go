package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
)

type UserAuthenRepository interface {
	GenericRepository[domain_entities.UserAuthen, domain_entities.UserAuthenFilter]
}
