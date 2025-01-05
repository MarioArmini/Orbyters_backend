package users

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	Id        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Users     []User `gorm:"many2many:user_roles;"`
}

func (r *Role) CreateRole(db *gorm.DB) error {
	return db.Create(&r).Error
}

func GetAllRoles(db *gorm.DB) ([]Role, error) {
	var roles []Role
	err := db.Model(&Role{}).Find(&roles).Error
	return roles, err
}

func GetRoleByName(db *gorm.DB, roleName string) (*Role, error) {
	var role Role
	err := db.Where("name = ?", roleName).First(&role).Error
	return &role, err
}
