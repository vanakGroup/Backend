package entities

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	CityID       uint   `gorm:"not null;index"`
	City         City   `gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	PostalCode   string `gorm:"type:varchar(50);not null"`
	UnitNumber   uint   `gorm:"not null"`
	ExactAddress string `gorm:"type:text;not null"`

	ParkingID uint     `gorm:"not null;uniqueIndex"`
	Parking   *Parking `gorm:"foreignKey:ParkingID;constraint:OnDelete:CASCADE"`
}
