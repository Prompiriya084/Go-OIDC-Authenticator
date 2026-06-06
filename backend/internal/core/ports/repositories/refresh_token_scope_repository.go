package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"context"
)

type RefreshTokenScopeRepository interface {
	GenericRepository[domain_entities.RefreshTokenScope, domain_entities.RefreshTokenScopeFilter]
	GetAllWithMaster(ctx context.Context, filter *domain_entities.RefreshTokenScopeFilter) ([]*domain_entities.RefreshTokenScope, error)
}
