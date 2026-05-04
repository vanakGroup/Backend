package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName     string `gorm:"type:varchar(50);not null"`
	LastName      string `gorm:"type:varchar(50);not null"`
	Phone         string `gorm:"type:varchar(15);uniqueIndex"`
	NationalCode  string `gorm:"type:varchar(20);uniqueIndex"`
	ShabaNumber   string `gorm:"type:varchar(26)"`
	PhoneVerified bool   `gorm:"default:false"`

	Roles []Role `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;"`

	Contacts []*User `gorm:"many2many:user_contacts;constraint:OnDelete:CASCADE;"`
}
