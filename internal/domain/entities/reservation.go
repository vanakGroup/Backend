package entities

import (
	"github.com/Vanak/Backend/internal/domain/enums"
	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	TotalPrice int64                   `gorm:"type:bigint; check:total_price >= 0; not null"`
	BasePrice  int64                   `gorm:"type:bigint;not null"`
	StartTime  int64                   `gorm:"not null"`
	EndTime    int64                   `gorm:"not null"`
	Hours      int                     `gorm:"not null"`
	Status     enums.ReservationStatus `gorm:"not null"`
	UserID     uint                    `gorm:"not null;index"`
	ParkingID  uint                    `gorm:"not null;index"`
	VehicleID  uint                    `gorm:"not null;index"`
	User       User                    `gorm:"foreignKey:UserID"`
	Parking    Parking                 `gorm:"foreignKey:ParkingID"`
	Vehicle    Vehicle                 `gorm:"foreignKey:VehicleID"`
	Slots      []AvailabilitySlot      `gorm:"foreignKey:ReservationID"`
}
