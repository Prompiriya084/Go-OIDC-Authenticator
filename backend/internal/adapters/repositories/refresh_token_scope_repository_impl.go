package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	"context"

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

func (r *refreshTokenScopeRepositoryImpl) GetAllWithMaster(
	ctx context.Context,
	filter *domain_entities.RefreshTokenScopeFilter,
) ([]*domain_entities.RefreshTokenScope, error) {
	var entities []*domain_entities.RefreshTokenScope
	query := r.db.WithContext(ctx).
		Preload("Token").
		Preload("Scope")
	if filter != nil {
		query = query.Where(filter)
	}
	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}
