package users

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id                 uint `gorm:"primaryKey"`
	Name               string
	Surname            string
	Email              *string `gorm:"unique;not null"`
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
	PassWordHash       string `gorm:"not null"`
	Roles              []Role `gorm:"many2many:user_roles;"`
	Reset_token        string
	Reset_token_expiry *time.Time
}

func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(&u).Error
}

func (u *User) GetUserByEmail(db *gorm.DB) (*User, error) {
	var user User
	err := db.Where("email = ?", u.Email).First(&user).Error
	return &user, err
}

func GetUserById(db *gorm.DB, userId uint) (*User, error) {
	var user User
	err := db.Where("Id = ?", userId).First(&user).Error
	return &user, err
}

func (u *User) HasRole(db *gorm.DB, roleId uint) (bool, error) {
	err := db.Preload("Roles").Find(u).Error
	if err != nil {
		return false, err
	}

	for _, role := range u.Roles {
		if role.Id == roleId {
			return true, nil
		}
	}

	return false, nil
}
