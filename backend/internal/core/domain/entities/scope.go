package domain_entities

import (
	"github.com/google/uuid"
)

// Scope represents the Master.Scopes table
type Scope struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AudienceID  *uuid.UUID `gorm:"type:uuid" json:"audienceId,omitempty"` // Nullable UUID
	Name        string     `gorm:"not null" json:"name"`
	ScopeTypeID uuid.UUID  `gorm:"not null;type:uuid" json:"scopeTypeId"`
}

func (Scope) TableName() string {
	return "Master.Scopes"
}

type ScopeFilter struct {
	ID          *uuid.UUID
	AudienceID  *uuid.UUID
	Name        *string
	ScopeTypeID *uuid.UUID
}
