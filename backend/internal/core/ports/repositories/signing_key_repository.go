package ports_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"context"
	"time"
)

type SigningKeyRepository interface {
	GenericRepository[domain_entities.SigningKey, domain_entities.SigningKeyFilter]
	GetActiveKey(ctx context.Context) (*domain_entities.SigningKey, error)
	GetAllKeys(ctx context.Context) ([]*domain_entities.SigningKey, error)
	GetExpiredKeys(ctx context.Context, cutoff time.Time) ([]*domain_entities.SigningKey, error)
}
