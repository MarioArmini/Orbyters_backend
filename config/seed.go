package config

import (
	"log"

	models "Orbyters/models/users"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "Admin"},
		{Name: "Customer"},
	}

	for _, role := range roles {
		var existingRole models.Role
		err := db.Where("name = ?", role.Name).First(&existingRole).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Fatalf("Error creating role: %v", err)
		}

		if existingRole.Id == 0 {
			existingRole.Name = role.Name
			if err := existingRole.CreateRole(db); err != nil {
				log.Fatalf("Error saving role: %v", err)
			}
			log.Printf("Role '%s' created.", role.Name)
		} else {
			log.Printf("Role '%s' already existing.", role.Name)
		}
	}
}
