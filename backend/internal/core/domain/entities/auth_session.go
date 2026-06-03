package domain_entities

import (
	"time"

	"github.com/google/uuid"
)

// AuthSession represents user active sessions
type AuthSession struct {
	SessionID   uuid.UUID `gorm:"primaryKey;type:uuid" json:"sessionId"`
	UserID      uuid.UUID `gorm:"not null;type:uuid" json:"userId"`
	CreatedAt   time.Time `gorm:"not null" json:"createdAt"`
	ExpiresAt   time.Time `gorm:"not null" json:"expiresAt"`
	CreatedAtTH time.Time `gorm:"not null;column:created_at_th" json:"createdAtTh"`
	ExpiresAtTH time.Time `gorm:"not null;column:expires_at_th" json:"expiresAtTh"`
}

type AuthSessionFilter struct {
	SessionID   *uuid.UUID
	UserID      *uuid.UUID
	CreatedAt   *time.Time
	ExpiresAt   *time.Time
	CreatedAtTH *time.Time
	ExpiresAtTH *time.Time
}
