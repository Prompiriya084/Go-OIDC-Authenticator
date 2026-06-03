package services

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"context"

	"github.com/google/uuid"
)

type MfaService interface {
	GetClientById(ctx context.Context, clientId uuid.UUID) (*domain_entities.Client, error)
	StartSetup(ctx context.Context, userId uuid.UUID) (string, error)
	ConfirmTotp(ctx context.Context, userId uuid.UUID, code string) (uuid.UUID, error)
	VerifyTotp(ctx context.Context, userId uuid.UUID, code string) (uuid.UUID, error)
}
