package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type viewRefreshTokenScopeRepositoryImpl struct {
	db *gorm.DB
}

func NewViewRefreshTokenScopeRepository(db *gorm.DB) ports_repositories.ViewRefreshTokenScopeRepository {
	return &viewRefreshTokenScopeRepositoryImpl{
		db: db,
	}
}
func (r *viewRefreshTokenScopeRepositoryImpl) GetAllByIDs(ctx context.Context, tokenID uuid.UUID) ([]domain_entities.ViewRefreshTokenScope, error) {
	var results []domain_entities.ViewRefreshTokenScope

	// GORM จะแกะฟังก์ชัน TableName() จาก Struct อัตโนมัติ
	// แล้วเปลี่ยนคำสั่งคิวรีให้กลายเป็น: SELECT * FROM View_RefreshToken_Scopes WHERE refresh_token_id = '...'
	err := r.db.WithContext(ctx).
		Where("token_id = ?", tokenID).
		Find(&results).
		Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
