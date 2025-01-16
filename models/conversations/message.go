package conversations

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	Id             uint   `gorm:"primaryKey"`
	Content        string `gorm:"not null"`
	CreatedAt      *time.Time
	ConversationId uint
	MessageTypeId  uint
	MessageType    MessageType `gorm:"foreignKey:MessageTypeId"`
	Role           string      `gorm:"not null"`
}

func (m *Message) CreteMessage(db *gorm.DB) error {
	return db.Create(&m).Error
}

func GetMessageById(db *gorm.DB, Id uint) (*Message, error) {
	var message Message
	err := db.Where("Id = ?", Id).Preload("MessageType").First(&message).Error
	return &message, err
}

func GetConversationHistory(db *gorm.DB, conversationId uint) ([]Message, error) {
	var messages []Message
	err := db.Where("conversation_id = ?", conversationId).Preload("MessageType").Find(&messages).Error
	return messages, err
}
