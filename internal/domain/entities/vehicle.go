package entities

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	LicensePlate        string `gorm:"type:varchar(20);uniqueIndex:idx_active_license_plate,where:deleted_at IS NULL"`
	CarModel            string `gorm:"type:varchar(50)"`
	RegistrationCardUrl string `gorm:"type:varchar(255)"`
	IsVerified          bool   `gorm:"default:false"`
	UserId              uint   `gorm:"not null;index"`
}
