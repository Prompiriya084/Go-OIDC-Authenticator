package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type userAuthenRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.UserAuthen, domain_entities.UserAuthenFilter]
}

func NewUserAuthenRepository(db *gorm.DB) ports_repositories.UserAuthenRepository {
	return &userAuthenRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.UserAuthen, domain_entities.UserAuthenFilter](db),
	}
}
