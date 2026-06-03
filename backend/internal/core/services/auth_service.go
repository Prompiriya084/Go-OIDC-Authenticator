package services

import (
	"OIDCAuthenticator/internal/core/dto"
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	ValidateGrantType(ctx context.Context, grantType string) bool
	CreateAuthorizationCode(ctx context.Context, sessionId uuid.UUID, state dto.AuthState) (string, error)
	HandleTokenAuthorizationCode(ctx context.Context, authorizationCode, clientId string, clientSecret *string, redirectUri, codeVerifier string) (*dto.TokenResult, error)
	HandleTokenRefreshToken(ctx context.Context, refreshToken, clientId string, clientSecret *string) (*dto.TokenResult, error)
}
