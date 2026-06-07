package services

import (
	"OIDCAuthenticator/internal/core/dto"
	"context"
)

type AccountService interface {
	SignIn(ctx context.Context, req dto.SignInRequestDTO) (*dto.SignInResponseDTO, error)
	GeneratePreMfaToken(userID string) (string, error)
	GenerateMfaToken(userID string) (string, error)
}
