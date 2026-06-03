package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"context"

	"github.com/google/uuid"
)

type ScopeRepository interface {
	GenericRepository[domain_entities.Scope, domain_entities.ScopeFilter]
	GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain_entities.Scope, error)
	GetAllByNames(ctx context.Context, names []string) ([]*domain_entities.Scope, error)
}
