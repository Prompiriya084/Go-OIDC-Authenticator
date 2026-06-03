package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type grantTypeRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.GrantType, domain_entities.GrantTypeFilter]
}

func NewGrantTypeRepository(db *gorm.DB) ports_repositories.GrantTypeRepository {
	return &grantTypeRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.GrantType, domain_entities.GrantTypeFilter](db),
	}
}
