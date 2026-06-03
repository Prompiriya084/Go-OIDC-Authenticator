package domain_entities

import (
	"github.com/google/uuid"
)

// RefreshTokenScope represents the dbo.RefreshToken_Scopes mapping table
type RefreshTokenScope struct {
	ID      int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TokenID uuid.UUID `gorm:"not null;type:uuid" json:"tokenId"`
	ScopeID uuid.UUID `gorm:"not null;type:uuid" json:"scopeId"`
}

func (RefreshTokenScope) TableName() string {
	return "dbo.RefreshToken_Scopes"
}

type RefreshTokenScopeFilter struct {
	ID      *int
	TokenID *uuid.UUID
	ScopeID *uuid.UUID
}
