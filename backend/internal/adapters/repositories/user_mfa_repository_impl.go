package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type userMfaRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.UserMfa, domain_entities.UserMfaFilter]
}

func NewUserMfaRepository(db *gorm.DB) ports_repositories.UserMfaRepository {
	return &userMfaRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.UserMfa, domain_entities.UserMfaFilter](db),
	}
}
