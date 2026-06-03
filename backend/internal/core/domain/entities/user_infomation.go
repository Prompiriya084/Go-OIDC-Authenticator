package domain_entities

import "github.com/google/uuid"

// UserInformation: ตารางเก็บข้อมูลส่วนตัวพนักงาน แยก Scope ออกจาก Auth ชัดเจน
type UserInformation struct {
	ID        uuid.UUID `gorm:"column:Id;primaryKey;type:uniqueidentifier" json:"id"`
	PrefixTH  string    `gorm:"column:Prefix_TH" json:"prefixTh"`
	PrefixEN  string    `gorm:"column:Prefix_EN" json:"prefixEn"`
	NameTH    string    `gorm:"column:Name_TH" json:"nameTh"`
	SurnameTH string    `gorm:"column:Surname_TH" json:"surnameTh"`
	NameEN    string    `gorm:"column:Name_EN" json:"nameEn"`
	SurnameEN string    `gorm:"column:Surname_EN" json:"surnameEn"`
	AddressTH string    `gorm:"column:Address_TH" json:"addressTH"`
	AddressEN string    `gorm:"column:Address_EN" json:"addressEN"`
}

func (UserInformation) TableName() string {
	return "dbo.User_Information"
}

type UserInformationFilter struct {
	ID        *uuid.UUID
	PrefixTH  *string
	PrefixEN  *string
	NameTH    *string
	SurnameTH *string
	NameEN    *string
	SurnameEN *string
	AddressTH *string
	AddressEN *string
}
