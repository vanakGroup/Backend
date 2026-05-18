package entities

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);index;not null"`
	Description string `gorm:"type:varchar(100)"`

	Roles []Role `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE"`
}
