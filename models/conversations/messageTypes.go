package conversations

import (
	"gorm.io/gorm"
)

type MessageType struct {
	Id   uint `gorm:"primaryKey"`
	Type string
}

func (msgType *MessageType) CreateMessageType(db *gorm.DB) error {
	return db.Create(&msgType).Error
}

func GetMessageTypeByType(db *gorm.DB, msgType string) (*MessageType, error) {
	var messageType MessageType
	err := db.Where("Type = ?", msgType).First(&messageType).Error
	return &messageType, err
}

func GetAllMessageTypes(db *gorm.DB) ([]MessageType, error) {
	var messageTypes []MessageType
	err := db.Model(&MessageType{}).Find(&messageTypes).Error
	return messageTypes, err
}
