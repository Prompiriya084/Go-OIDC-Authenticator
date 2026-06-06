package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"context"
)

type ClientGrantTypeRepository interface {
	GenericRepository[domain_entities.ClientGrantType, domain_entities.ClientGrantTypeFilter]
	GetAllWithMaster(ctx context.Context, filter *domain_entities.ClientGrantTypeFilter) ([]*domain_entities.ClientGrantType, error)
}
