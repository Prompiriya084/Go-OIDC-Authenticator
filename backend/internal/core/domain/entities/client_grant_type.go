package domain_entities

import (
	"github.com/google/uuid"
)

// ClientGrantType represents the dbo.Client_GrantTypes mapping table
type ClientGrantType struct {
	ID       int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID uuid.UUID `gorm:"not null;type:uuid" json:"clientId"`
	GrantID  uuid.UUID `gorm:"not null;type:uuid" json:"grantId"`

	Client Client    `gorm:"foreignKey:ClientID;references:ID" json:"client,omitempty"`
	Grant  GrantType `gorm:"foreignKey:GrantID;references:ID" json:"grant,omitempty"`
}

func (ClientGrantType) TableName() string {
	return "dbo.Client_GrantTypes"
}

type ClientGrantTypeFilter struct {
	ID       *int
	ClientID *uuid.UUID `gorm:"column:client_id"`
	GrantID  *uuid.UUID `gorm:"column:grant_id"`
}
