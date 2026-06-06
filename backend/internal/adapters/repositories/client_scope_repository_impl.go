package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	"context"

	"gorm.io/gorm"
)

type clientScopeRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.ClientScope, domain_entities.ClientScopeFilter]
}

func NewClientScopeRepository(db *gorm.DB) ports_repositories.ClientScopeRepository {
	return &clientScopeRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.ClientScope, domain_entities.ClientScopeFilter](db),
	}
}

func (r *clientScopeRepositoryImpl) GetAllWithMaster(
	ctx context.Context,
	filter *domain_entities.ClientScopeFilter,
) ([]*domain_entities.ClientScope, error) {
	var entities []*domain_entities.ClientScope
	query := r.db.WithContext(ctx).
		Preload("Client").
		Preload("Scope")
	if filter != nil {
		query = query.Where(filter)
	}
	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}
