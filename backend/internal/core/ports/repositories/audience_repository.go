package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"context"

	"github.com/google/uuid"
)

type AudienceRepository interface {
	GenericRepository[domain_entities.Audience, domain_entities.AudienceFilter]
	GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]domain_entities.Audience, error)
}
