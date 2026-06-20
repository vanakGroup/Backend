package entities

import "gorm.io/gorm"

type MockPayment struct {
	gorm.Model
	Amount int64 `gorm:"type:bigint;default:0;check:amount >= 0;not null"`
	UserID uint  `gorm:"not null;index"`
}
