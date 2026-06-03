package domain_entities

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents the dbo.RefreshTokens table
type RefreshToken struct {
	ID                  uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	SessionID           uuid.UUID `gorm:"not null;type:uuid" json:"sessionId"`
	UserID              uuid.UUID `gorm:"not null;type:uuid" json:"userId"`
	ClientID            uuid.UUID `gorm:"not null;type:uuid" json:"clientId"`
	TokenHash           string    `gorm:"not null" json:"tokenHash"`
	CreatedAt           time.Time `gorm:"not null" json:"createdAt"`
	ExpiresAt           time.Time `gorm:"not null" json:"expiresAt"`
	InitialSignInDate   time.Time `gorm:"not null" json:"initialSignInDate"`
	CreatedAtTH         time.Time `gorm:"not null;column:created_at_th" json:"createdAtTh"`
	ExpiresAtTH         time.Time `gorm:"not null;column:expires_at_th" json:"expiresAtTh"`
	InitialSignInDateTH time.Time `gorm:"not null;column:initial_sign_in_date_th" json:"initialSignInDateTh"`
	IsRevoked           bool      `gorm:"not null" json:"isRevoked"`
}

func (RefreshToken) TableName() string {
	return "dbo.RefreshTokens"
}

// RefreshToken represents the dbo.RefreshTokens table
type RefreshTokenFilter struct {
	ID                  *uuid.UUID
	SessionID           *uuid.UUID
	UserID              *uuid.UUID
	ClientID            *uuid.UUID
	TokenHash           *string
	CreatedAt           *time.Time
	ExpiresAt           *time.Time
	InitialSignInDate   *time.Time
	CreatedAtTH         *time.Time
	ExpiresAtTH         *time.Time
	InitialSignInDateTH *time.Time
	IsRevoked           *bool
}
