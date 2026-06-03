package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type refreshTokenScopeRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.RefreshTokenScope, domain_entities.RefreshTokenScopeFilter]
}

func NewRefreshTokenScopeRepository(db *gorm.DB) ports_repositories.RefreshTokenScopeRepository {
	return &refreshTokenScopeRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.RefreshTokenScope, domain_entities.RefreshTokenScopeFilter](db),
	}
}
