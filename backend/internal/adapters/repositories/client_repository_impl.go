package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type clientRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.Client, domain_entities.ClientFilter]
}

func NewClientRepository(db *gorm.DB) ports_repositories.ClientRepository {
	return &clientRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.Client, domain_entities.ClientFilter](db),
	}
}
