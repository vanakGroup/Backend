package entities

import "gorm.io/gorm"

type City struct {
	gorm.Model
	Name       string `gorm:"type:varchar(50);not null"`
	ProvinceID uint   `gorm:"not null"`
}
