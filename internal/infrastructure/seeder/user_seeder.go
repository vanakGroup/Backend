// this file not use because need the RBAC then we should fix that

package seeder

import (
	"log"

	"github.com/Vanak/Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type UserSeeder struct {
	db *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

func (us *UserSeeder) SeedRoles() {
	roles := []entities.Role{
		{Name: "customer"},
		{Name: "admin"},
	}

	for _, role := range roles {
		var existingRole entities.Role
		result := us.db.Where("type = ?", role.Name).First(&existingRole)

		if result.Error == gorm.ErrRecordNotFound {
			if err := us.db.Create(&role).Error; err != nil {
				log.Printf("Error creating role %s: %v", role.Name, err)
			} else {
				log.Printf("Role '%s' created successfully", role.Name)
			}
		} else if result.Error != nil {
			log.Printf("Error checking role %s: %v", role.Name, result.Error)
		}
	}
}
