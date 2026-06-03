package dto

import domain_entities "OIDCAuthenticator/internal/core/domain/entities"

type RefreshTokenResult struct {
	PlainTextToken string
	Entity         domain_entities.RefreshToken
}
