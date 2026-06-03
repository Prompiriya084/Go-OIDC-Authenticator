package domain_entities

import (
	"github.com/google/uuid"
)

// Client represents the Master.Clients table
type Client struct {
	ID                          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name                        string    `gorm:"not null" json:"name"`
	HashSecret                  *string   `json:"hashSecret,omitempty"`
	RedirectURI                 string    `gorm:"not null;column:redirect_uri" json:"redirectUri"`
	DefaultRedirectURI          string    `gorm:"not null;column:default_redirect_uri" json:"defaultRedirectUri"`
	RequiredClientSecret        bool      `gorm:"not null" json:"requiredClientSecret"`
	RequirePCKE                 bool      `gorm:"not null;column:require_pcke" json:"requirePcke"` // แก้อักษรย่อ PKCE ตามโค้ดต้นฉบับ
	IsActive                    bool      `gorm:"not null" json:"isActive"`
	AccessTokenLifeTimeMinutes  int       `gorm:"not null" json:"accessTokenLifeTimeMinutes"`
	RefreshTokenLifeTimeMinutes int       `gorm:"not null" json:"refreshTokenLifeTimeMinutes"`
}

func (Client) TableName() string {
	return "Master.Clients"
}

type ClientFilter struct {
	ID                          *uuid.UUID
	Name                        *string
	HashSecret                  *string
	RedirectURI                 *string
	DefaultRedirectURI          *string
	RequiredClientSecret        *bool
	RequirePCKE                 *bool
	IsActive                    *bool
	AccessTokenLifeTimeMinutes  *int
	RefreshTokenLifeTimeMinutes *int
}
