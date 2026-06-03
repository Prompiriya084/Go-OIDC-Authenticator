package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type scopeRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.Scope, domain_entities.ScopeFilter]
}

func NewScopeRepository(db *gorm.DB) ports_repositories.ScopeRepository {
	return &scopeRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.Scope, domain_entities.ScopeFilter](db),
	}
}
func (r *scopeRepositoryImpl) GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain_entities.Scope, error) {
	var scopes []*domain_entities.Scope

	// ถ้าไม่มี IDs ส่งมาเลย ให้รีเทิร์นอาเรย์ว่างกลับไปทันที ไม่ต้องยิงคิวรีให้เปลืองทรัพยากร
	if len(ids) == 0 {
		return scopes, nil
	}

	// GORM จะฉลาดพอที่จะแปลงคำสั่ง .Where("id IN ?", ids) ให้เป็น SQL "SELECT * FROM scopes WHERE id IN (...)"
	// พร้อมแนบ ctx ไปทำ Request Cancellation / Timeout ให้อัตโนมัติ
	err := r.db.WithContext(ctx).
		Where("ID IN ?", ids).
		Find(&scopes).
		Error

	if err != nil {
		return nil, err
	}

	return scopes, nil
}
func (r *scopeRepositoryImpl) GetAllByNames(ctx context.Context, names []string) ([]*domain_entities.Scope, error) {
	var scopes []*domain_entities.Scope

	// ถ้าไม่มี IDs ส่งมาเลย ให้รีเทิร์นอาเรย์ว่างกลับไปทันที ไม่ต้องยิงคิวรีให้เปลืองทรัพยากร
	if len(names) == 0 {
		return scopes, nil
	}

	// GORM จะฉลาดพอที่จะแปลงคำสั่ง .Where("id IN ?", ids) ให้เป็น SQL "SELECT * FROM scopes WHERE id IN (...)"
	// พร้อมแนบ ctx ไปทำ Request Cancellation / Timeout ให้อัตโนมัติ
	err := r.db.WithContext(ctx).
		Where("Name IN ?", names).
		Find(&scopes).
		Error

	if err != nil {
		return nil, err
	}

	return scopes, nil
}
