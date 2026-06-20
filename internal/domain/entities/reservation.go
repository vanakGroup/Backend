package entities

import (
	"github.com/Vanak/Backend/internal/domain/enums"
	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	TotalPrice int64                   `gorm:"type:bigint; check:total_price >= 0; not null"`
	StartTime  uint                    `gorm:"not null"`
	EndTime    uint                    `gorm:"not null"`
	Status     enums.ReservationStatus `gorm:"not null"`
	UserID     uint                    `gorm:"not null;index"`
	ParkingID  uint                    `gorm:"not null;index"`
	VehicleID  uint                    `gorm:"not null;index"`
}
