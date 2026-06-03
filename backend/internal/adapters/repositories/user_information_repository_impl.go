package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type userInformationRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.UserInformation, domain_entities.UserInformationFilter]
}

func NewUserInformationRepository(db *gorm.DB) ports_repositories.UserInformationRepository {
	return &userInformationRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.UserInformation, domain_entities.UserInformationFilter](db),
	}
}
