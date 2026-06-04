package services

import (
	"OIDCAuthenticator/internal/core/dto"
	"context"
)

type AuthService interface {
	ValidateGrantType(ctx context.Context, grantType string) bool
	Authorize(ctx context.Context, req dto.AuthorizeRequestDTO, flowID string, sessionID string) (*dto.AuthorizeResult, error)
	HandleTokenAuthorizationCode(ctx context.Context, authorizationCode, clientId string, clientSecret *string, redirectUri, codeVerifier string) (*dto.TokenResult, error)
	HandleTokenRefreshToken(ctx context.Context, refreshToken, clientId string, clientSecret *string) (*dto.TokenResult, error)
}
