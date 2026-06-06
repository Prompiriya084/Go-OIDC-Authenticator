package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	"context"

	"gorm.io/gorm"
)

type clientGrantTypeRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.ClientGrantType, domain_entities.ClientGrantTypeFilter]
}

func NewClientGrantTypeRepository(db *gorm.DB) ports_repositories.ClientGrantTypeRepository {
	return &clientGrantTypeRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.ClientGrantType, domain_entities.ClientGrantTypeFilter](db),
	}
}

func (r *clientGrantTypeRepositoryImpl) GetAllWithMaster(
	ctx context.Context,
	filter *domain_entities.ClientGrantTypeFilter,
) ([]*domain_entities.ClientGrantType, error) {
	var entities []*domain_entities.ClientGrantType
	query := r.db.WithContext(ctx).
		Preload("Client").
		Preload("Grant")
	if filter != nil {
		query = query.Where(filter)
	}
	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}
