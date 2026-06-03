package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type clientScopeRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.ClientScope, domain_entities.ClientScopeFilter]
}

func NewClientScopeRepository(db *gorm.DB) ports_repositories.ClientScopeRepository {
	return &clientScopeRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.ClientScope, domain_entities.ClientScopeFilter](db),
	}
}
