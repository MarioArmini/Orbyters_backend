package config

import (
	"log"

	conversationsModels "Orbyters/models/conversations"
	userModels "Orbyters/models/users"

	"gorm.io/gorm"
)

func ApplySeeds(db *gorm.DB) {
	seedRoles(db)
	seedMessageTypes(db)
}

func seedRoles(db *gorm.DB) {
	roles := []userModels.Role{
		{Name: "Admin"},
		{Name: "Customer"},
	}

	for _, role := range roles {
		var existingRole userModels.Role
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

func seedMessageTypes(db *gorm.DB) {
	messageTypes := []conversationsModels.MessageType{
		{Type: "user"},
		{Type: "assistant"},
	}

	for _, messageType := range messageTypes {
		var existingMessageType conversationsModels.MessageType
		err := db.Where("Type = ?", messageType.Type).First(&existingMessageType).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Fatalf("Error creating message type: %v", err)
		}

		if existingMessageType.Id == 0 {
			existingMessageType.Type = messageType.Type
			if err := existingMessageType.CreateMessageType(db); err != nil {
				log.Fatalf("Error saving message type: %v", err)
			}
			log.Printf("Message type '%s' created.", messageType.Type)
		} else {
			log.Printf("Message type '%s' already existing.", messageType.Type)
		}
	}
}
