package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

type refreshTokenRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.RefreshToken, domain_entities.RefreshTokenFilter]
}

func NewRefreshTokenRepository(db *gorm.DB) ports_repositories.RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.RefreshToken, domain_entities.RefreshTokenFilter](db),
	}
}
