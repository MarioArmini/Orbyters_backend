package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id           uint `gorm:"primaryKey"`
	Name         string
	Surname      string
	Email        *string `gorm:"unique;not null"`
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	PassWordHash string `gorm:"not null"`
}

func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(&u).Error
}

func (u *User) GetUserByEmail(db *gorm.DB) (*User, error) {
	var user User
	err := db.Where("email = ?", u.Email).First(&user).Error
	return &user, err
}
