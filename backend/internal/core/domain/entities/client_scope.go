package domain_entities

import (
	"github.com/google/uuid"
)

// ClientScope represents the dbo.Client_Scopes mapping table
type ClientScope struct {
	ID       int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID uuid.UUID `gorm:"not null;type:uuid" json:"clientId"`
	ScopeID  uuid.UUID `gorm:"not null;type:uuid" json:"scopeId"`

	Client Client `gorm:"foreignKey:ClientID;references:ID" json:"client,omitempty"`
	Scope  Scope  `gorm:"foreignKey:ScopeID;references:ID" json:"scope,omitempty"`
}

func (ClientScope) TableName() string {
	return "dbo.Client_Scopes"
}

type ClientScopeFilter struct {
	ID       *int
	ClientID *uuid.UUID
	ScopeID  *uuid.UUID
}
