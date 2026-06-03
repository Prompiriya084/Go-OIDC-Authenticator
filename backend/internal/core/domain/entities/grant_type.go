package domain_entities

import (
	"github.com/google/uuid"
)

// GrantType represents the Master.GrantTypes table
type GrantType struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Type        string    `gorm:"not null" json:"type"`
	Description *string   `json:"description,omitempty"`
}

func (GrantType) TableName() string {
	return "Master.GrantTypes"
}

type GrantTypeFilter struct {
	ID          *uuid.UUID
	Type        *string
	Description *string
}
