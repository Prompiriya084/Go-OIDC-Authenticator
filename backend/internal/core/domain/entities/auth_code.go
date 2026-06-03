package domain_entities

import (
	"time"

	"github.com/google/uuid"
)

// AuthCode represents the auth code session data
type AuthCode struct {
	Code            string    `gorm:"primaryKey" json:"code"`
	SessionID       uuid.UUID `gorm:"not null;type:uuid" json:"sessionId"`
	UserID          uuid.UUID `gorm:"not null;type:uuid" json:"userId"`
	ClientID        uuid.UUID `gorm:"not null;type:uuid" json:"clientId"`
	CodeChallenge   string    `gorm:"not null" json:"codeChallenge"`
	ChallengeMethod string    `gorm:"not null" json:"challengeMethod"`
	RequiredScopes  string    `gorm:"not null" json:"requiredScopes"`
	RedirectURI     *string   `json:"redirectUri,omitempty"`
	Nonce           *string   `json:"nonce,omitempty"`
	ExpiresAt       time.Time `gorm:"not null" json:"expiresAt"`
	ExpiresAtTH     time.Time `gorm:"not null;column:expires_at_th" json:"expiresAtTh"`
}

// AuthCode represents the auth code session data
type AuthCodeFilter struct {
	Code            *string
	SessionID       *uuid.UUID
	UserID          *uuid.UUID
	ClientID        *uuid.UUID
	CodeChallenge   *string
	ChallengeMethod *string
	RequiredScopes  *string
	RedirectURI     *string
	Nonce           *string
	ExpiresAt       *time.Time
	ExpiresAtTH     *time.Time
}
