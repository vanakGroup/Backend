package entities

import "gorm.io/gorm"

type Province struct {
	gorm.Model
	Name   string `gorm:"type:varchar(50);not null"`
	Cities []City `gorm:"foreignKey:ProvinceID"`
}
