package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type authSessionRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.AuthSession, domain_entities.AuthSessionFilter]
}

func NewAuthSessionRepository(db *gorm.DB) ports_repositories.AuthSessionRepository {
	return &authSessionRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.AuthSession, domain_entities.AuthSessionFilter](db),
	}
}
