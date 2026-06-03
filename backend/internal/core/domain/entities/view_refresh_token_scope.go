package domain_entities

import "github.com/google/uuid"

type ViewRefreshTokenScope struct {
	TokenID   uuid.UUID `gorm:"column:token_id"`
	ScopeID   uuid.UUID `gorm:"column:scope_id"`
	ScopeName uuid.UUID `gorm:"column:scope_name"`
	ScopeType uuid.UUID `gorm:"column:scope_type"`
}

func (ViewRefreshTokenScope) TableName() string {
	return "View_RefreshToken_Scopes"
}
