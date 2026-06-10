package entities

import (
	"time"

	"gorm.io/gorm"
)

type AvailabilitySlot struct {
	gorm.Model
	Date     time.Time `gorm:"not null"`
	Hour     uint      `gorm:"not null"`
	IsBooked bool      `gorm:"default:false"`

	ParkingID uint    `gorm:"not null"`
	Parking   Parking `gorm:"foreignKey:ParkingID;constraint:OnDelete:CASCADE"`
}
