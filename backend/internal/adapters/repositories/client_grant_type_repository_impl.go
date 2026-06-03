package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type clientGrantTypeImpl struct {
	*genericRepositoryImpl[domain_entities.ClientGrantType, domain_entities.ClientGrantTypeFilter]
}

func NewClientGrantTypeRepository(db *gorm.DB) ports_repositories.ClientGrantTypeRepository {
	return &clientGrantTypeImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.ClientGrantType, domain_entities.ClientGrantTypeFilter](db),
	}
}
