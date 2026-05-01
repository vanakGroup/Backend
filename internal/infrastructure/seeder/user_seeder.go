package seeder

import (
	"gorm.io/gorm"
)

type UserSeeder struct {
	db *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

// func (userSeeder *UserSeeder) SeedRoles() {
// 	roles := []entities.Role{
// 		{Type: "customer"},
// 		{Type: "admin"},
// 		{Type: "moderator"},
// 	}

// 	for _, role := range roles {
// 		var existingRole entities.Role
// 		result := userSeeder.db.Where("type = ?", role.Type).First(&existingRole)

// 		if result.Error == gorm.ErrRecordNotFound {
// 			if err := userSeeder.db.Create(&role).Error; err != nil {
// 				log.Printf("Error creating role %s: %v", role.Type, err)
// 			} else {
// 				log.Printf("Role '%s' created successfully", role.Type)
// 			}
// 		} else if result.Error != nil {
// 			log.Printf("Error checking role %s: %v", role.Type, result.Error)
// 		}
// 	}
// }
