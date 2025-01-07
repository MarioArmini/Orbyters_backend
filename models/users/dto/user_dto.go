package dto

import (
	models "Orbyters/models/users"
)

type UserDto struct {
	Id      uint `gorm:"primaryKey"`
	Name    string
	Surname string
	Email   *string       `gorm:"unique;not null"`
	Roles   []models.Role `gorm:"many2many:user_roles;"`
}
