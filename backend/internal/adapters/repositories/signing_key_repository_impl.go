package adapters_repositories

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	"context"
	"time"

	"gorm.io/gorm"
)

type signingKeyRepositoryImpl struct {
	*genericRepositoryImpl[domain_entities.SigningKey, domain_entities.SigningKeyFilter]
}

func NewSigningKeyRepository(db *gorm.DB) ports_repositories.SigningKeyRepository {
	return &signingKeyRepositoryImpl{
		genericRepositoryImpl: newGenericRepositoryImpl[domain_entities.SigningKey, domain_entities.SigningKeyFilter](db),
	}
}

// 1. GetActiveKey โดยใช้ Filter Struct (แทนการเขียน Hard Query "is_active = ?")
func (r *signingKeyRepositoryImpl) GetActiveKey(ctx context.Context) (*domain_entities.SigningKey, error) {
	// ประกาศ Filter และส่งเฉพาะค่า IsActive: true เข้าไป (ฟิลด์อื่นจะเป็น nil)
	activeTrue := true
	filter := &domain_entities.SigningKeyFilter{
		IsActive: &activeTrue,
	}

	var entity domain_entities.SigningKey

	// เนื่องจากต้องการ Order และดึงเอาตัวล่าสุดตัวเดียว เราจึงครอบด้วย db context ปกติ
	// แต่เงื่อนไข .Where() เราส่งตัว filter struct เข้าไปตรงๆ ได้เลย!
	err := r.db.WithContext(ctx).
		Where(filter).
		Order("created_at DESC").
		First(&entity).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

// 2. GetAllKeys ดึงทั้งหมดเหมือนเดิม
func (r *signingKeyRepositoryImpl) GetAllKeys(ctx context.Context) ([]*domain_entities.SigningKey, error) {
	return r.GetAll(ctx, nil)
}

// 3. GetExpiredKeys โดยใช้ Filter Struct ร่วมกับเงื่อนไขพิเศษ
func (r *signingKeyRepositoryImpl) GetExpiredKeys(ctx context.Context, cutoff time.Time) ([]*domain_entities.SigningKey, error) {
	activeFalse := false
	filter := &domain_entities.SigningKeyFilter{
		IsActive: &activeFalse, // ใช้ Struct คุมเงื่อนไข IsActive = false
	}

	var entities []*domain_entities.SigningKey

	// สำหรับเรื่องเวลา (CreatedAt < cutoff) ซึ่งเป็นเงื่อนไขแบบเปรียบเทียบ (Less Than)
	// เราจะเอามาเขียนเสริมต่อท้ายได้แบบนี้ ซึ่งอ่านง่ายและปลอดภัยเหมือนเดิมครับ
	err := r.db.WithContext(ctx).
		Where(filter).
		Where("created_at < ?", cutoff).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}
