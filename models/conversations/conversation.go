package conversations

import (
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	Id        uint `gorm:"primaryKey"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	UserId    uint
	Messages  []Message
}

func (c *Conversation) CreateConversation(db *gorm.DB) error {
	return db.Create(&c).Error
}

func GetConversationById(db *gorm.DB, Id uint) (*Conversation, error) {
	var conversation Conversation
	err := db.Where("Id = ?", Id).Preload("Messages").First(&conversation).Error
	return &conversation, err
}

func GetConversationsByUserId(db *gorm.DB, userId uint) ([]Conversation, error) {
	var conversations []Conversation
	err := db.Where("UserId = ?", userId).Preload("Messages").Find(&conversations).Error
	return conversations, err
}

func (c *Conversation) AppendMessage(db *gorm.DB, message *Message) []Message {
	var conversation Conversation
	newMessages := append(conversation.Messages, *message)
	return newMessages
}
