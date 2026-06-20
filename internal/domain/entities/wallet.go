package entities

import (
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Balance      int64         `gorm:"type:bigint;default:0;check:balance >= 0;not null"`
	IsActive     bool          `gorm:"default:false; not null"`
	UserID       uint          `gorm:"not null;uniqueindex"`
	Transactions []Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
