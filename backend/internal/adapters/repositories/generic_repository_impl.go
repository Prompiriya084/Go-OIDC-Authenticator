package adapters_repositories

import (
	"context"

	"gorm.io/gorm"
)

type genericRepositoryImpl[Tentity any, Tfilter any] struct {
	db *gorm.DB
}

func newGenericRepositoryImpl[Tentity any, Tfilter any](db *gorm.DB) *genericRepositoryImpl[Tentity, Tfilter] {
	return &genericRepositoryImpl[Tentity, Tfilter]{db: db}
}

func (r *genericRepositoryImpl[Tentity, Tfilter]) Add(ctx context.Context, entity *Tentity) error {
	return r.db.WithContext(ctx).Create(entity).Error // 🔹 แก้บั๊ก & ออกแล้ว
}
func (r *genericRepositoryImpl[Tentity, Tfilter]) AddRange(ctx context.Context, entities []*Tentity) error {
	// ถ้าส่งสไลซ์ว่างมา ไม่ต้องวิ่งไปยิงคิวรีในฐานข้อมูลให้เสียเวลา
	if len(entities) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&entities).Error
}

func (r *genericRepositoryImpl[Tentity, Tfilter]) Update(ctx context.Context, entity *Tentity) error {
	return r.db.WithContext(ctx).Save(entity).Error // 🔹 แก้บั๊ก & ออกแล้ว
}
func (r *genericRepositoryImpl[Tentity, Tfilter]) UpdateRange(ctx context.Context, entities []*Tentity) error {
	if len(entities) == 0 {
		return nil
	}

	// ⚠️ ระวัง: คำสั่ง Save() ใน GORM จะทำงานแบบ Upsert (ถ้าไม่มี ID จะสร้างใหม่ ถ้ามี ID จะอัปเดต)
	// และมันจะยิงอัปเดตแยกทีละแถว (Single Queries) ไม่ใช่ Batch Update ยัดรอบเดียวแบบ Create!
	for _, entity := range entities {
		if entity == nil {
			continue
		}
		if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *genericRepositoryImpl[Tentity, Tfilter]) Delete(ctx context.Context, entity *Tentity) error {
	return r.db.WithContext(ctx).Delete(entity).Error // 🔹 แก้บั๊ก & ออกแล้ว
}
func (r *genericRepositoryImpl[Tentity, Tfilter]) DeleteRange(ctx context.Context, entities []*Tentity) error {
	if len(entities) == 0 {
		return nil
	}

	// GORM จะแกะเอา Primary Key ของทุก Object ใน Slice มารวมร่าง
	// แล้วยิงคำสั่งเป็น "DELETE FROM table WHERE id IN (1, 2, 3)" รอบเดียวจบ!
	return r.db.WithContext(ctx).Delete(&entities).Error
}

func (r *genericRepositoryImpl[Tentity, Tfilter]) Get(ctx context.Context, filters *Tfilter) (*Tentity, error) {
	var entity Tentity

	// 🔹 แก้ปัญหา Data Leak ด้วยการสร้าง Scope ใหม่ผ่าน WithContext
	query := r.db.WithContext(ctx)
	if filters != nil {
		query = query.Where(filters)
	}

	if err := query.First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *genericRepositoryImpl[Tentity, Tfilter]) GetAll(ctx context.Context, filters *Tfilter) ([]*Tentity, error) {
	var entities []*Tentity

	query := r.db.WithContext(ctx)
	if filters != nil {
		query = query.Where(filters)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}
