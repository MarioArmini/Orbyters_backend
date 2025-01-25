package users

import (
	"time"

	"gorm.io/gorm"
)

type UserSubscription struct {
	Id             uint `gorm:"primaryKey"`
	UserId         uint
	SubscriptionId uint
	ExpiresAt      *time.Time
}

func (us *UserSubscription) CreateUserSubscription(db *gorm.DB) error {
	return db.Create(&us).Error
}

func HasSubscription(db *gorm.DB, userId uint, subScriptionId uint) (bool, error) {
	var userSubscription UserSubscription
	err := db.Where(&UserSubscription{SubscriptionId: subScriptionId, UserId: userId}).First(&userSubscription).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
