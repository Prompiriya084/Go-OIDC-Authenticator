package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type audienceRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.Audience, domain_entities.AudienceFilter]
}

func NewAudienceRepository(db *gorm.DB) ports_repositories.AudienceRepository {
	return &audienceRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.Audience, domain_entities.AudienceFilter](db),
	}
}
func (r *audienceRepositoryImpl) GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]domain_entities.Audience, error) {
	var audiences []domain_entities.Audience

	// ถ้าไม่มี IDs ส่งมาเลย ให้รีเทิร์นอาเรย์ว่างกลับไปทันที ไม่ต้องยิงคิวรีให้เปลืองทรัพยากร
	if len(ids) == 0 {
		return audiences, nil
	}

	// GORM จะฉลาดพอที่จะแปลงคำสั่ง .Where("id IN ?", ids) ให้เป็น SQL "SELECT * FROM scopes WHERE id IN (...)"
	// พร้อมแนบ ctx ไปทำ Request Cancellation / Timeout ให้อัตโนมัติ
	err := r.db.WithContext(ctx).
		Where("ID IN ?", ids).
		Find(&audiences).
		Error

	if err != nil {
		return nil, err
	}

	return audiences, nil
}
