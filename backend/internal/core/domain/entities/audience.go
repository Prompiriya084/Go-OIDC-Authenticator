package domain_entities

import (
	"github.com/google/uuid"
)

// Audience represents the Master.Audiences table
type Audience struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description *string   `json:"description,omitempty"` // ใช้ pointer สำหรับ nullable string
}

// TableName overrides the table name for GORM
func (Audience) TableName() string {
	return "Master.Audiences"
}

type AudienceFilter struct {
	ID          *uuid.UUID
	Name        *string
	Description *string
}
