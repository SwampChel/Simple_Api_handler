package userService

import (
	"gorm.io/gorm"
	"pet_project_1_etap/internal/taskService"
)

type User struct {
	gorm.Model
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Tasks    []taskService.Task `json:"tasks"`
}
