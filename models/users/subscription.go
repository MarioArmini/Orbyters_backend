package users

import (
	"gorm.io/gorm"
)

type Subscription struct {
	Id              uint    `gorm:"primaryKey" json:"id"`
	Price           float64 `json:"price"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Users           []User  `gorm:"many2many:user_subscriptions;"`
	StripeProductId string  `json:"stripeProductId"`
}

func (s *Subscription) CreateSubscription(db *gorm.DB) error {
	return db.Create(&s).Error
}

func GetSubscriptionById(db *gorm.DB, Id uint) (*Subscription, error) {
	var subscription Subscription
	err := db.Where("Id = ?", Id).First(&subscription).Error
	return &subscription, err
}

func GetAllSubscriptions(db *gorm.DB) ([]Subscription, error) {
	var subscriptions []Subscription
	err := db.Model(&Subscription{}).Find(&subscriptions).Error
	return subscriptions, err
}
