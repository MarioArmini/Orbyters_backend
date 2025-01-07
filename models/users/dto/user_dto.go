package dto

import (
	models "Orbyters/models/users"
)

type UserDto struct {
	Id      uint          `json:"id"`
	Name    string        `json:"name"`
	Surname string        `json:"surname"`
	Email   *string       `json:"email"`
	Roles   []models.Role `json:"roles"`
}
