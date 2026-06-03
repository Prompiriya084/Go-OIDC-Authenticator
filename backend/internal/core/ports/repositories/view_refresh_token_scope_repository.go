package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"context"

	"github.com/google/uuid"
)

type ViewRefreshTokenScopeRepository interface {
	GetAllByIDs(ctx context.Context, tokenID uuid.UUID) ([]domain_entities.ViewRefreshTokenScope, error)
}
