package entities

import (
	"github.com/Vanak/Backend/internal/domain/enums"
	"gorm.io/gorm"
)

type Parking struct {
	gorm.Model
	Title       string              `gorm:"size:255;not null"`
	Latitude    float64             `gorm:"type:decimal(10,8);not null"`
	Longitude   float64             `gorm:"type:decimal(11,8);not null"`
	BasePrice   float64             `gorm:"type:decimal(10,2);not null"`
	Status      enums.ParkingStatus `gorm:"default:2"`
	DocumentURL string              `gorm:"type:varchar(500);not null"`

	UserID       uint               `gorm:"not null;index"`
	User         User               `gorm:"foreignKey:UserID"`
	Address      Address            `gorm:"foreignKey:ParkingID;constraint:OnDelete:CASCADE"`
	Attributes   []ParkingAttribute `gorm:"many2many:parking_attribute_relations;constraint:OnDelete:CASCADE"`
	Slots        []AvailabilitySlot `gorm:"foreignKey:ParkingID;constraint:OnDelete:CASCADE"`
	Images       []ParkingImage     `gorm:"foreignKey:ParkingID;constraint:OnDelete:CASCADE"`
	Reservations []Reservation      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
