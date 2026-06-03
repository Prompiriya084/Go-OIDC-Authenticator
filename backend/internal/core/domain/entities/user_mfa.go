package domain_entities

import (
	"time"

	"github.com/google/uuid"
)

type UserMfa struct {
	// UserId: ตัวนี้ทำหน้าที่เป็น PK ของตารางตัวเอง **แต่ห้ามทำ autoIncrement เด็ดขาด** // เพราะค่าของมันจะต้องถูกส่งมาจาก UserId ของ UserAuthen (ตารางหลัก) เสมอตอนสมัครใช้งาน
	ID                  uuid.UUID  `gorm:"column:UserId;primaryKey;type:uniqueidentifier;autoIncrement:false" json:"userId"`
	TotpSecretEncrypted string     `gorm:"column:TotpSecretEncrypted" json:"totpSecretEncrypted"`
	TotpEnabled         bool       `gorm:"column:TotpEnabled;not null" json:"totpEnabled"`
	TotpConfirmedAt     *time.Time `gorm:"column:TotpConfirmedAt" json:"totpConfirmedAt"`
	LastMfaAt           *time.Time `gorm:"column:LastMfaAt" json:"lastMfaAt"`
	TotpConfirmedAtTH   *time.Time `gorm:"column:TotpConfirmedAt_TH" json:"totpConfirmedAtTh"`
	LastMfaAtTH         *time.Time `gorm:"column:LastMfaAt_TH" json:"lastMfaAtTh"`
}

func (UserMfa) TableName() string {
	return "dbo.User_Mfa"
}

type UserMfaFilter struct {
	// UserId: ตัวนี้ทำหน้าที่เป็น PK ของตารางตัวเอง **แต่ห้ามทำ autoIncrement เด็ดขาด** // เพราะค่าของมันจะต้องถูกส่งมาจาก UserId ของ UserAuthen (ตารางหลัก) เสมอตอนสมัครใช้งาน
	ID                  *uuid.UUID
	TotpSecretEncrypted *string
	TotpEnabled         *bool
	TotpConfirmedAt     *time.Time
	LastMfaAt           *time.Time
	TotpConfirmedAtTH   *time.Time
	LastMfaAtTH         *time.Time
}
