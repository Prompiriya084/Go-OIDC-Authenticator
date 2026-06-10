package services

import (
	"OIDCAuthenticator/internal/core/dto"
	"context"

	"github.com/google/uuid"
)

type MfaService interface {
	GetOIDCFlowState(ctx context.Context, flowId string) (*dto.OIDCFlowState, error)
	GetDefaultURIByClientId(ctx context.Context, clientId uuid.UUID) (string, error)
	StartSetup(ctx context.Context, userId uuid.UUID) (string, error)
	ConfirmTotp(ctx context.Context, userId uuid.UUID, code string) (*dto.MfaResponseDTO, error)
	VerifyTotp(ctx context.Context, userId uuid.UUID, code string) (*dto.MfaResponseDTO, error)
}
