package domain_entities

import (
	"time"

	"github.com/google/uuid"
)

// UserAuthen: ตารางหลักสำหรับเก็บข้อมูล Account เพื่อใช้ Log in
type UserAuthen struct {
	// UserId: แนะนำให้ใช้ uuid.UUID และ "ไม่ใช่" Auto Increment ของ DB (เพราะเป็น UUID)
	// แต่ใส่ default:newid() ไว้ให้ GORM รู้ว่าถ้าส่งค่าว่างมา ให้ใช้ฟังก์ชันของ DB เจนค่าให้
	ID           uuid.UUID `gorm:"column:Id;primaryKey;type:uniqueidentifier;default:newid()" json:"id"`
	Username     string    `gorm:"column:Username;not null" json:"username"`
	PasswordHash string    `gorm:"column:PasswordHash" json:"passwordHash"`
	SignedInAt   time.Time `gorm:"column:SignedInAt" json:"signedInAt"`
	ExpiresAt    time.Time `gorm:"column:ExpiresAt" json:"expiresAt"`
	SignedInAtTH time.Time `gorm:"column:SignedInAt_TH" json:"signedInAtTh"`
	ExpiresAtTH  time.Time `gorm:"column:ExpiresAt_TH" json:"expiresAtTh"`
	IsActive     bool      `gorm:"column:IsActive" json:"isActive"`
}

func (UserAuthen) TableName() string {
	return "dbo.User_Authen"
}

type UserAuthenFilter struct {
	ID           *uuid.UUID
	Username     *string
	PasswordHash *string
	SignedInAt   *time.Time
	ExpiresAt    *time.Time
	SignedInAtTH *time.Time
	ExpiresAtTH  *time.Time
	IsActive     *bool
}
