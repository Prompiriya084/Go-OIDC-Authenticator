package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"context"
)

type ClientScopeRepository interface {
	GenericRepository[domain_entities.ClientScope, domain_entities.ClientScopeFilter]
	GetAllWithMaster(ctx context.Context, filter *domain_entities.ClientScopeFilter) ([]*domain_entities.ClientScope, error)
}
