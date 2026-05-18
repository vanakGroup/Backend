package entities

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);uniqueIndex;not null"`
	Description string `gorm:"type:varchar(100)"`

	Users       []User       `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE"`
}
