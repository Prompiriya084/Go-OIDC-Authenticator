package services

import (
	"OIDCAuthenticator/internal/core/dto"
	"context"
)

type AuthService interface {
	ValidateGrantType(ctx context.Context, grantType string) bool
	Authorize(ctx context.Context, req dto.AuthorizeRequestDTO, flowID string, sessionID string) (*dto.AuthorizeResult, error)
	HandleToken(ctx context.Context, req dto.TokenRequestDTO) (*dto.TokenResponseDTO, error)
}
