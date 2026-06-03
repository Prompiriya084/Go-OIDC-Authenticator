package domain_entities

import "time"

type SigningKey struct {
	Id         string `gorm:"primaryKey"`
	Kid        string `gorm:"index"`
	PrivateKey string // เก็บเป็น Base64 string
	PublicKey  string // เก็บเป็น Base64 string
	CreatedAt  time.Time
	IsActive   bool
}

type SigningKeyFilter struct {
	Id         *string
	Kid        *string
	PrivateKey *string
	PublicKey  *string
	CreatedAt  *time.Time
	IsActive   *bool
}
