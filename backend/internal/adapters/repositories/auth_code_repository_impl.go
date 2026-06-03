package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"

	"gorm.io/gorm"
)

// --- ADAPTERS ---
type authCodeRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.AuthCode, domain_entities.AuthCodeFilter]
}

func NewAuthCodeRepository(db *gorm.DB) ports_repositories.AuthCodeRepository {
	return &authCodeRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.AuthCode, domain_entities.AuthCodeFilter](db),
	}
}
