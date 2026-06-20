package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName    *string `gorm:"type:varchar(50)"`
	LastName     *string `gorm:"type:varchar(50)"`
	Phone        string  `gorm:"type:varchar(15);uniqueIndex;not null"`
	NationalCode *string `gorm:"type:varchar(20);uniqueIndex"`
	ShabaNumber  *string `gorm:"type:varchar(26);uniqueIndex"`

	PhoneVerified    bool `gorm:"default:false"`
	MinInfoCompleted bool `gorm:"default:false"`

	Roles        []Role        `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;"`
	Contacts     []*User       `gorm:"many2many:user_contacts;constraint:OnDelete:CASCADE;"`
	Parkings     []Parking     `gorm:"foreignKey:UserID"`
	Vehicles     []Vehicle     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Wallet       Wallet        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payments     []MockPayment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reservations []Reservation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
