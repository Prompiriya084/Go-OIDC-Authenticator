package dataaccess

import (
	ports_database "OIDCAuthenticator/internal/core/ports/database"
	"context"

	"gorm.io/gorm"
)

type gormTxKey struct{}

type txManagerImpl struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) ports_database.TransactionManager {
	return &txManagerImpl{db: db}
}

func (m *txManagerImpl) Begin(ctx context.Context) (context.Context, error) {
	tx := m.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	// 💡 แอบฝัง tx ไว้ใน Context ตัวใหม่
	return context.WithValue(ctx, gormTxKey{}, tx), nil
}

func (m *txManagerImpl) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(gormTxKey{}).(*gorm.DB); ok {
		return tx.Commit().Error
	}
	return nil
}

func (m *txManagerImpl) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value(gormTxKey{}).(*gorm.DB); ok {
		return tx.Rollback().Error
	}
	return nil
}

// 💡 ฟังก์ชันช่วยสำหรับฝั่ง Repository ดึง DB ออกไปใช้งาน
func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(gormTxKey{}).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}
	return defaultDB.WithContext(ctx)
}
