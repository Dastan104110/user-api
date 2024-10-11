package models

import (
	"gorm.io/gorm"
)

// User - модель пользователя
type User struct {
	gorm.Model        // Встраиваемая модель GORM
	Name       string `json:"name" gorm:"unique;not null"` // Имя пользователя
	Age        int    `json:"age"`                         // Возраст пользователя
}
