package entities

import "gorm.io/gorm"

type ParkingImage struct {
	gorm.Model
	ImageURL  string  `gorm:"type:varchar(500);not null"`
	ParkingID uint    `gorm:"not null;index"`
	Parking   Parking `gorm:"foreignKey:ParkingID;constraint:OnDelete:CASCADE"`
}
