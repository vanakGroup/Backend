package entities

import (
	"time"

	"gorm.io/gorm"
)

type AvailabilitySlot struct {
	gorm.Model
	StartTime time.Time `gorm:"not null;index"`
	EndTime   time.Time `gorm:"not null;index"`
	IsBooked  bool      `gorm:"default:false;index"`

	ParkingID uint    `gorm:"not null"`
	Parking   Parking `gorm:"foreignKey:ParkingID;constraint:OnDelete:CASCADE"`
}
