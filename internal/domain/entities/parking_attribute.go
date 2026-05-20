package entities

import "gorm.io/gorm"

type ParkingAttribute struct {
	gorm.Model
	Name string `gorm:"type:varchar(20)"`

	Parkings []Parking `gorm:"many2many:parking_attributes;constraint:OnDelete:CASCADE;"`
}
